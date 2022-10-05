package config

import (
	"github.com/minelytix/config-api/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database   DatabaseConfig   `yaml:"database"`
	Collection CollectionConfig `yaml:"collection"`
	Logging    log.Config       `yaml:"logging"`
	Server     ServerConfig     `yaml:"server"`
	Health     HealthConfig     `yaml:"health"`
}

type DatabaseConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type CollectionConfig struct {
	config string `yaml:"config"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type HealthConfig struct {
	Port int `yaml:"port"`
}

func LoadConfig(data []byte) (*Config, error) {
	var conf = Config{}
	var err = yaml.Unmarshal(data, &conf)
	return &conf, err
}
