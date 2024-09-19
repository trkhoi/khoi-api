package config

import (
	"log"

	"github.com/spf13/viper"
)

// Loader load config from reader into Viper
type Loader interface {
	Load(*viper.Viper) (*viper.Viper, error)
}

// Read sets default to a viper instance and read user config to override these defaults
func Read() (*viper.Viper, error) {
	dcfg := viper.New()

	// read configs from .env file
	fileLoader := NewFileLoader(".env", ".")
	dcfg, err := fileLoader.Load(dcfg)
	if err != nil {
		log.Printf("Failed to load .env file: %s", err.Error())
	}

	// read configs from environment variables
	envLoader := NewENVLoader()
	dcfg, err = envLoader.Load(dcfg)
	if err != nil {
		log.Printf("Failed to load environment: %s", err.Error())
	}

	return dcfg, nil
}
