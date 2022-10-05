package config

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/middleware"
	"github.com/minelytix/config-api/api"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigApi interface {
	ListConfigs(pagination api.Pagination) (*ConfigDtoCollection, error)
	GetConfig(string) (*ConfigDto, error)
	CreateConfig(*ConfigDto) (*ConfigDto, error)
	UpdateConfig(string, *ConfigDto) (*ConfigDto, error)
	DeleteConfig(string) error
}

type ConfigApiAdapter func(ConfigApi) ConfigApi

func AdaptConfigApi(r ConfigApi, adapters ...ConfigApiAdapter) ConfigApi {
	for _, adapter := range adapters {
		r = adapter(r)
	}
	return r
}

type ConfigService struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewConfigService(ctx context.Context, collection *mongo.Collection) ConfigApi {
	getReqId := middleware.GetReqID(ctx)

	return &ConfigService{
		ctx:        context.WithValue(ctx, middleware.RequestIDKey, getReqId),
		collection: collection,
	}
}

func (s *ConfigService) ListConfigs(pagination api.Pagination) (*ConfigDtoCollection, error) {
	configDao := NewMongoDbConfigDao(s.ctx, s.collection)

	models, err := configDao.List(pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, err
	}
	count, err := configDao.Count()
	if err != nil {
		return nil, err
	}

	dtos, err := ModelsToDtos(models)
	if err != nil {
		return nil, err
	}

	result := ConfigDtoCollection{
		CollectionMeta: api.CollectionMeta{
			Total: *count,
		},
		ConfigDtos: dtos,
	}

	return &result, err
}

func (s *ConfigService) GetConfig(id string) (*ConfigDto, error) {
	configDao := NewMongoDbConfigDao(s.ctx, s.collection)
	model, err := configDao.Get(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &api.InfoNotFound{Entity: fmt.Sprintf("config with id = %s", id)}
		}
		return nil, err
	}

	return ModelToDto(model)
}

func (s *ConfigService) CreateConfig(dto *ConfigDto) (*ConfigDto, error) {
	configDao := NewMongoDbConfigDao(s.ctx, s.collection)

	if configDao.Exists(dto.GetId()) {
		return nil, &api.WarningValidation{Msg: fmt.Sprintf("config with id = %s already exists", dto.GetId())}
	}

	payload, err := DtoToModel(dto)
	if err != nil {
		return nil, err
	}

	model, err := configDao.Create(payload)
	if err != nil {
		return nil, err
	}

	return ModelToDto(model)
}

func (s *ConfigService) UpdateConfig(id string, dto *ConfigDto) (*ConfigDto, error) {
	configDao := NewMongoDbConfigDao(s.ctx, s.collection)

	if _, err := s.GetConfig(id); err != nil {
		return nil, err
	}

	payload, err := DtoAndIdToModel(dto, id)
	if err != nil {
		return nil, err
	}

	model, err := configDao.Update(id, payload)
	if err != nil {
		return nil, err
	}

	return ModelToDto(model)
}

func (s *ConfigService) DeleteConfig(id string) error {
	configDao := NewMongoDbConfigDao(s.ctx, s.collection)
	return configDao.Delete(id)
}
