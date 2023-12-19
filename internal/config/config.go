package config

import (
	"go-server/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"` // Обязательно нужно указывать
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`         // Есть дефолт значения
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"` // Есть дефолт значения
		Port   string `yaml:"port" env-default:"true"`         // Есть дефолт значения
	} `yaml:"listen"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config { // SINGLETON
	once.Do(func() { // once.Do сработает 1 раз, при последующих вызовах - сразу вернет instance
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}

	})
	return instance
}
