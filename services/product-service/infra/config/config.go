package config

type AppConfig struct {
	Port string `yaml:"port" mapstructure:"port"`
}
