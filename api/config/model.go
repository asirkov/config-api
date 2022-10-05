package config

import (
	"encoding/json"
	"fmt"

	"github.com/minelytix/config-api/api"
)

type ConfigDtoCollection struct {
	api.CollectionMeta

	ConfigDtos []*ConfigDto `json:"configs"`
}

type ConfigDto map[string]interface{}

func (d ConfigDto) GetId() string {
	return d["$id"].(string)
}

func (d ConfigDto) SetId(id string) {
	d["$id"] = id
}

func (d ConfigDto) GetConfig() map[string]interface{} {
	bytes, err := json.Marshal(d)
	if err != nil {
		return nil
	}

	var config map[string]interface{}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil
	}

	return config
}

// func (d ConfigDto) GetConfigVersion() string {
// 	return d["$config"].(string)
// }

type ConfigModel map[string]interface{}

func (m ConfigModel) GetId() string {
	return m["_id"].(string)
}

func (m ConfigModel) GetConfig() map[string]interface{} {
	bytes, err := json.Marshal(m["config"])
	if err != nil {
		return nil
	}

	var config map[string]interface{}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil
	}

	return config
}

func DtoToModel(dto *ConfigDto) (*ConfigModel, error) {
	if dto == nil {
		return nil, fmt.Errorf("failed to map config dto to model, model is nil")
	}

	return DtoAndIdToModel(dto, dto.GetId())
}

func DtoAndIdToModel(dto *ConfigDto, id string) (*ConfigModel, error) {
	if dto == nil {
		return nil, fmt.Errorf("failed to map config dto to model, model is nil")
	}
	dto.SetId(id)

	model := ConfigModel{
		"_id":    id,
		"config": dto.GetConfig(),
	}

	return &model, nil
}

func ModelToDto(model *ConfigModel) (*ConfigDto, error) {
	if model == nil {
		return nil, fmt.Errorf("failed to map config model to dto, dto is nil")
	}

	dto := ConfigDto(model.GetConfig())
	return &dto, nil
}

func ModelsToDtos(models []*ConfigModel) ([]*ConfigDto, error) {
	if models == nil {
		return nil, fmt.Errorf("failed to map config model to dtos, dto sis nil")
	}

	var dtos []*ConfigDto
	for _, model := range models {
		dto, err := ModelToDto(model)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}
