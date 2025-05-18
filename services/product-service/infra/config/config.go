package config

type AppConfig struct {
	Port             string `yaml:"port" mapstructure:"port"`
	REDIS_ADDR       string `yaml:"redis_addr" mapstructure:"redis_addr"`
	REDIS_PASS       string `yaml:"redis_pass" mapstructure:"redis_pass"`
	MONGO_URI        string `yaml:"mongo_uri" mapstructure:"mongo_uri"`
	MONGO_DB         string `yaml:"mongo_db" mapstructure:"mongo_db"`
	MONGO_COLLECTION string `yaml:"mongo_collection" mapstructure:"mongo_collection"`
}
