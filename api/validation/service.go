package validation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/minelytix/config-api/api/config"
	"github.com/santhosh-tekuri/jsonconfig/v5"
)

type ValidationApi interface {
	ValidateConfig(string, *map[string]interface{}) (*map[string]interface{}, error)
}

type ValidationApiAdapter func(ValidationApi) ValidationApi

func AdaptValidationApi(r ValidationApi, adapters ...ValidationApiAdapter) ValidationApi {
	for _, adapter := range adapters {
		r = adapter(r)
	}
	return r
}

type ValidationService struct {
	ctx       context.Context
	configApi config.ConfigApi
}

func NewValidationService(ctx context.Context, configApi config.ConfigApi) ValidationApi {
	return &ValidationService{
		ctx:       ctx,
		configApi: configApi,
	}
}

func (s *ValidationService) ValidateConfig(id string, document *map[string]interface{}) (*map[string]interface{}, error) {
	config, err := s.configApi.GetConfig(id)
	if err != nil {
		return nil, err
	}

	configBytes, err := json.Marshal(config.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to marshall config data from request body: %s", err)
	}

	dataBytes, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall document data from request body: %s", err)
	}

	var dataObj interface{}
	if err := json.Unmarshal(dataBytes, &dataObj); err != nil {
		return nil, fmt.Errorf("failed to marshall data from request body: %s", err)
	}

	configObj, err := jsonconfig.CompileString("", string(configBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create jsonconfig: %s", err)
	}

	if err := configObj.Validate(dataObj); err != nil {
		return nil, err
	}

	return document, nil
}
