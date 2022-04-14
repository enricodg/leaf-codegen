package config

import (
	leafConfig "github.com/paulusrobin/leaf-utilities/config"
	"log"
	"sync"
)

var (
	configuration EnvConfig
	once          sync.Once
)

type (
	EnvConfig struct {
		// - Logger config
		LogFilePath  string `envconfig:"LOG_FILE_NAME"`
		LogFormatter string `envconfig:"LOG_FORMATTER" default:"TEXT"`
	}
)

func GetConfig() EnvConfig {
	once.Do(func() {
		if err := leafConfig.NewFromEnv(&configuration); err != nil {
			log.Fatal(err)
		}
	})
	return configuration
}
