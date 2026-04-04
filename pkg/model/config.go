package model

type LogConfig struct {
	Level  string `mapstructure:"level" json:"level" yaml:"level"`
	Format string `mapstructure:"format" json:"format" yaml:"format"`
}

type ServerConfig struct {
	Port int `mapstructure:"port" json:"port" yaml:"port"`
}

type MongoConfig struct {
	URI    string `mapstructure:"uri" json:"uri" yaml:"uri"`
	DBName string `mapstructure:"db" json:"db" yaml:"db"`
}

type OtelConfig struct {
	Endpoint    string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	ServiceName string `mapstructure:"service_name" json:"service_name" yaml:"service_name"`
}
