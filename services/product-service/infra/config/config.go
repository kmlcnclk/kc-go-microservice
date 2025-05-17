package config

type AppConfig struct {
	Port       string `yaml:"port" mapstructure:"port"`
	REDIS_ADDR string `yaml:"redis_addr" mapstructure:"redis_addr"`
	REDIS_PASS string `yaml:"redis_pass" mapstructure:"redis_pass"`
}
