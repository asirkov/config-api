package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	configapi "github.com/minelytix/config-api/api/config"
	healthapi "github.com/minelytix/config-api/api/health"
	staticapi "github.com/minelytix/config-api/api/static"
	validationapi "github.com/minelytix/config-api/api/validation"
	"github.com/minelytix/config-api/log"
	"github.com/minelytix/config-api/router"

	"github.com/minelytix/config-api/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var wg sync.WaitGroup

func loadConfig(fileName string) (*config.Config, error) {
	if len(fileName) == 0 {
		return nil, fmt.Errorf("failed to load config: name of file is empty")
	}

	configData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %s", err)
	}

	conf, err := config.LoadConfig(configData)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %s", err)
	}

	return conf, nil
}

func initLogger(conf log.Config) error {
	loggers := []log.Logger{}

	if conf.LogFilesEnabled {
		zapFileLogger, err := log.NewZapFileLogger(conf)
		if err != nil {
			return fmt.Errorf("failed to initialize file logger: %s", err.Error())
		}
		loggers = append(loggers, zapFileLogger)
	}

	if conf.ConsoleEnabled {
		zapConsoleLogger, err := log.NewZapConsoleLogger(conf)
		if err != nil {
			return fmt.Errorf("failed to initialize console logger: %s", err.Error())
		}
		loggers = append(loggers, zapConsoleLogger)
	}

	multiLogAdapter := log.NewMultiLogAdapter(loggers...)
	log.Init(conf, multiLogAdapter)

	return nil
}

func initDb(conf *config.DatabaseConfig, ctx context.Context) (*mongo.Database, error) {
	options := options.Client().ApplyURI(conf.Uri)

	client, err := mongo.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongodb client: %s", err)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := client.Connect(timeoutCtx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err)
	}

	return client.Database(conf.Database), nil
}

func initConfigCollection(conf *config.CollectionConfig, db *mongo.Database) (*mongo.Collection, error) {
	configCollectionName := conf.config
	collection := db.Collection(configCollectionName)
	if collection == nil {
		return nil, fmt.Errorf("failed to get collection: %s", configCollectionName)
	}

	return collection, nil
}

func initHealthRouter(conf *config.HealthConfig, api healthapi.HealthApi) *router.Router {
	router := router.NewRouter(conf.Port)

	router.RegisterController(
		healthapi.NewHealthController(api),
	)

	return router
}

func initMainRouter(
	conf *config.ServerConfig,
	configApi configapi.ConfigApi,
	validationApi validationapi.ValidationApi,
) *router.Router {
	router := router.NewRouter(conf.Port)

	router.RegisterController(
		staticapi.NewStaticController(),
		configapi.NewConfigController(configApi),
		validationapi.NewValidationController(validationApi),
	)

	return router
}

func init() {
	_ = os.Setenv("TZ", "GMT")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSignals(sigs, &wg)
}

func handleSignals(c chan os.Signal, wg *sync.WaitGroup) {
	switch <-c {
	case syscall.SIGINT, syscall.SIGTERM:
		log.Info("Start to shutdown")
		wg.Wait()

		log.Info("Done shutdown")
	}

	os.Exit(0)
}

func main() {
	// Config
	conf, err := loadConfig("config/application.yml")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %s", err))
	}

	if err := initLogger(conf.Logging); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %s", err))
	}

	ctx := context.Background()

	db, err := initDb(&conf.Database, ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize config collection: %s", err))
	}
	defer db.Client().Disconnect(ctx)

	healthApi := healthapi.NewHealthService(db)

	healthRouter := initHealthRouter(&conf.Health, healthApi)
	go func() {
		log.Infof("Health router starts listening on port %d", conf.Health.Port)
		if err := healthRouter.Listen(); err != nil {
			log.Error("Health router listening error:", err)
		} else {
			log.Info("Health router listening is stopped")
		}
	}()

	configCollection, err := initConfigCollection(&conf.Collection, db)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize config collection: %s", err))
	}

	//******************************************************************************************************************
	// config
	//******************************************************************************************************************

	configValidator := configapi.ConfigValidator{}
	configService := configapi.NewConfigService(ctx, configCollection)
	configApi := configapi.AdaptConfigApi(configService, configValidator.ConfigApiAdapter)

	//******************************************************************************************************************
	// Validation
	//******************************************************************************************************************

	validationService := validationapi.NewValidationService(ctx, configApi)
	validationApi := validationapi.AdaptValidationApi(validationService)

	//*****************************************************************************************************************

	mainRouter := initMainRouter(&conf.Server, configApi, validationApi)

	log.Infof("Http router starts listening on port %d", conf.Server.Port)
	if err := mainRouter.Listen(); err != nil {
		log.Error("Http router listening error:", err)
	} else {
		log.Info("Http router listening is stopped")
	}

	log.Info("Shutting down gracefully...")

	if err := mainRouter.Shutdown(); err != nil {
		log.Error("Cannot shut down server gracefully. Forcing exit.", err)
	}

	log.Info("Server stopped")
}
