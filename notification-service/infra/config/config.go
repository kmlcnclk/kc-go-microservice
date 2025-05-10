package config

type AppConfig struct {
	RabbitMQURL      string `yaml:"rabbitmq_url" mapstructure:"rabbitmq_url"`
	RabbitOrderQueue string `yaml:"rabbitmq_order_queue" mapstructure:"rabbitmq_order_queue"`
}
