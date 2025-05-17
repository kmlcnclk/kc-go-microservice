package config

type AppConfig struct {
	Port                 string `yaml:"port" mapstructure:"port"`
	RabbitMQURL          string `yaml:"rabbitmq_url" mapstructure:"rabbitmq_url"`
	RabbitMQQueueName    string `yaml:"rabbitmq_queue_name" mapstructure:"rabbitmq_queue_name"`
	RabbitMQExchangeName string `yaml:"rabbitmq_exchange_name" mapstructure:"rabbitmq_exchange_name"`
	RabbitMQRoutingKey   string `yaml:"rabbitmq_routingKey" mapstructure:"rabbitmq_routingKey"`
}
