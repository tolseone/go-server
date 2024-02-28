package config

import (
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"go-server/pkg/logging"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"` // Обязательно нужно указывать
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`         // Есть дефолт значения
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"` // Есть дефолт значения
		Port   string `yaml:"port" env-default:"8080"`         // Есть дефолт значения
	} `yaml:"listen"`
	Storage   StorageConfig `yaml:"storage"`
	Clients   ClientsConfig `yaml:"clients"`
	AppSecret string        `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
}

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
}

type ClientsConfig struct {
	Auth Client `yaml:"auth"`
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
