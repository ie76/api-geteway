package external

import (
	"assignment/internal/errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type servicesMap struct {
	Services map[string]Config `yaml:"services"`
}

type Config struct {
	ServiceName     string                 `yaml:"aut_method"`
	BaseUrl         string                 `yaml:"base_url"`
	CacheDuration   int                    `yaml:"cache_duration"`
	AuthCredentials map[string]interface{} `yaml:"auth_credentials"`
}

func GetServiceConfig(name string) (*Config, *errors.Error) {
	absPath, err := filepath.Abs("./external/services.yaml")
	if err != nil {
		return nil, errors.New(errors.ErrConfigInit)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, errors.New(errors.ErrConfigInit)
	}

	var services servicesMap
	err = yaml.Unmarshal(data, &services)
	if err != nil {
		return nil, errors.New(errors.ErrConfigInit)
	}

	service, exists := services.Services[name]
	if !exists {
		return nil, errors.New(errors.ErrServiceNotFound)
	}

	config := Config{
		ServiceName:     service.ServiceName,
		BaseUrl:         service.BaseUrl,
		CacheDuration:   service.CacheDuration,
		AuthCredentials: service.AuthCredentials,
	}

	return &config, nil
}
