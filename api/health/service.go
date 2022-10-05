package health

import (
	"context"
	"fmt"

	"github.com/minelytix/config-api/api"
	"github.com/minelytix/config-api/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthApi interface {
	HealthCheck(ctx context.Context) ([]api.ServiceStatus, bool)
}

type HealthService struct {
	db *mongo.Database
}

func NewHealthService(db *mongo.Database) HealthApi {
	return &HealthService{
		db: db,
	}
}

func (r *HealthService) pingDb() error {
	return r.db.Client().Ping(context.TODO(), nil)
}

func (r *HealthService) HealthCheck(ctx context.Context) ([]api.ServiceStatus, bool) {
	statuses := []api.ServiceStatus{}
	failed := false

	msg := ""
	pFailed := false

	// Database
	if err := r.pingDb(); err != nil {
		failed = true
		pFailed = true
		msg = err.Error()
		log.ErrorCtx(ctx, fmt.Sprintf("Failed to Ping DataBase, %s", msg), nil)
	}
	status := api.ServiceStatus{Service: "Connection to database", Failed: pFailed, Message: msg}
	statuses = append(statuses, status)

	return statuses, failed
}
