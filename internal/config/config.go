package config

import (
	"mado/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// The Singleton design pattern ensures that a class has only one instance and provides a global point of access to that instance.
type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	MongoDB struct {
		Host       string `json:""`
		Port       string `json:""`
		Database   string `json:""`
		AuthDB     string `json:""`
		Username   string `json:""`
		Password   string `json:""`
		Collection string `json:""`
	}
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() { // Do is intended for initialization that must be run exactly once.
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
