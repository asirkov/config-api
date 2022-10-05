package config

import (
	"fmt"

	"github.com/minelytix/config-api/api"
)

type ConfigValidator struct {
	api ConfigApi
}

func (v *ConfigValidator) ConfigApiAdapter(api ConfigApi) ConfigApi {
	v.api = api
	return v
}

func (v *ConfigValidator) ListConfigs(pagination api.Pagination) (*ConfigDtoCollection, error) {
	return v.api.ListConfigs(pagination)
}

func (v *ConfigValidator) GetConfig(id string) (*ConfigDto, error) {
	return v.api.GetConfig(id)
}

func (v *ConfigValidator) CreateConfig(dto *ConfigDto) (*ConfigDto, error) {
	if err := validateConfigCreate(dto); err != nil {
		return nil, err
	}

	return v.api.CreateConfig(dto)
}

func (v *ConfigValidator) UpdateConfig(id string, dto *ConfigDto) (*ConfigDto, error) {
	if err := validateConfigUpdate(dto); err != nil {
		return nil, err
	}

	return v.api.UpdateConfig(id, dto)
}

func (v *ConfigValidator) DeleteConfig(id string) error {
	return v.api.DeleteConfig(id)
}

var (
	validConfigVersion = map[interface{}]bool{
		"http://json-config.org/draft-04/config#":      true,
		"http://json-config.org/draft-06/config#":      true,
		"http://json-config.org/draft-07/config#":      true,
		"https://json-config.org/draft/2019-09/config": true,
	}
)

func validateConfigCreate(dto *ConfigDto) error {
	if dto == nil {
		return &api.WarningValidation{Msg: "object is nil"}
	}
	value := *dto

	if value["$id"] == nil || value["$id"] == "" {
		return &api.WarningValidation{Msg: "field '$id' is missed"}
	}
	if value["$config"] == nil || value["$config"] == "" {
		return &api.WarningValidation{Msg: "field '$config' is missed"}
	}
	if value["$config"] != nil && !validConfigVersion[value["$config"]] {
		return &api.WarningValidation{Msg: fmt.Sprintf("the value %s is not allowed for field '$config'", value["$config"])}
	}

	return nil
}

func validateConfigUpdate(dto *ConfigDto) error {
	return validateConfigCreate(dto)
}
