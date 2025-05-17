package config

type AppConfig struct {
	Port                string `yaml:"port" mapstructure:"port"`
	OrderServiceAddress string `yaml:"order_service_address" mapstructure:"order_service_address"`
}
