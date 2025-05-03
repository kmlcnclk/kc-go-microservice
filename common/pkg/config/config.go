package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig[T any](configPath ...string) *T {
	var config T

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/config")

	for _, path := range configPath {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	return &config
}
