package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var globalConfig Configuration

func GetConfig() Configuration {
	return globalConfig
}

func LoadConfig(name string) error {
	cfg := Configuration{}

	if name != "" {
		viper.SetConfigName(name)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./internal/infrastructure/config")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				panic(fmt.Errorf("fatal error config file: %w", err))
			} else {
				panic(fmt.Errorf("fatal error config file: %w", err))
			}
		}

		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		globalConfig = cfg
	}

	return nil
}
