package config

import (
	"context"
	"github.com/bsonger/devflow-common/client/logging"
	"github.com/bsonger/devflow-common/client/mongo"
	commonModel "github.com/bsonger/devflow-common/model"
	"github.com/bsonger/devflow-service-common/observability"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/bsonger/devflow-verify-service/pkg/store"
	"github.com/spf13/viper"
)

type Config struct {
	Server    *model.ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
	Mongo     *model.MongoConfig  `mapstructure:"mongo"  json:"mongo"  yaml:"mongo"`
	Log       *model.LogConfig    `mapstructure:"log"    json:"log"    yaml:"log"`
	Otel      *model.OtelConfig   `mapstructure:"otel"   json:"otel"   yaml:"otel"`
	Pyroscope string              `mapstructure:"pyroscope" json:"pyroscope" yaml:"pyroscope"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")
	v.AddConfigPath("/etc/devflow/config/")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config *Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func InitConfig(ctx context.Context, config *Config) error {
	_, err := InitRuntime(ctx, config, "")
	return err
}

func InitRuntime(ctx context.Context, config *Config, serviceName string) (func(context.Context) error, error) {
	shutdown, err := observability.Init(ctx, observability.RuntimeOptions{
		LogLevel:        stringValue(config.Log, func(v *model.LogConfig) string { return v.Level }),
		LogFormat:       stringValue(config.Log, func(v *model.LogConfig) string { return v.Format }),
		OtelEndpoint:    stringValue(config.Otel, func(v *model.OtelConfig) string { return v.Endpoint }),
		OtelService:     stringValue(config.Otel, func(v *model.OtelConfig) string { return v.ServiceName }),
		PyroscopeAddr:   configValue(config, func(v *Config) string { return v.Pyroscope }),
		ServiceOverride: serviceName,
	})
	if err != nil {
		return nil, err
	}

	client, err := mongo.InitMongo(ctx, toCommonMongoConfig(config.Mongo), logging.Logger)
	if err != nil {
		return shutdown, err
	}
	store.InitMongo(client, config.Mongo.DBName)
	return shutdown, nil
}

func toCommonMongoConfig(cfg *model.MongoConfig) *commonModel.MongoConfig {
	if cfg == nil {
		return nil
	}
	return &commonModel.MongoConfig{
		URI:    cfg.URI,
		DBName: cfg.DBName,
	}
}

func ResolveConfigPort(cfg *Config) int {
	if cfg == nil || cfg.Server == nil {
		return 0
	}
	return cfg.Server.Port
}

func stringValue[T any](value *T, getter func(*T) string) string {
	if value == nil {
		return ""
	}
	return getter(value)
}

func configValue(cfg *Config, getter func(*Config) string) string {
	if cfg == nil {
		return ""
	}
	return getter(cfg)
}
