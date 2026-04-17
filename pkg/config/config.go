package config

import (
	"context"
	"database/sql"

	"github.com/bsonger/devflow-service-common/observability"
	"github.com/bsonger/devflow-verify-service/pkg/api"
	"github.com/bsonger/devflow-verify-service/pkg/model"
	"github.com/bsonger/devflow-verify-service/pkg/store"
	"github.com/spf13/viper"
)

type Config struct {
	Server    *model.ServerConfig   `mapstructure:"server" json:"server" yaml:"server"`
	Postgres  *model.PostgresConfig `mapstructure:"postgres" json:"postgres" yaml:"postgres"`
	Log       *model.LogConfig      `mapstructure:"log" json:"log" yaml:"log"`
	Otel      *model.OtelConfig     `mapstructure:"otel" json:"otel" yaml:"otel"`
	Auth      *model.AuthConfig     `mapstructure:"auth" json:"auth" yaml:"auth"`
	Pyroscope string                `mapstructure:"pyroscope" json:"pyroscope" yaml:"pyroscope"`
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

	db, err := sql.Open("pgx", stringValue(config.Postgres, func(v *model.PostgresConfig) string { return v.DSN }))
	if err != nil {
		return shutdown, err
	}
	store.ApplyPool(db,
		intValue(config.Postgres, func(v *model.PostgresConfig) int { return v.MaxOpenConns }),
		intValue(config.Postgres, func(v *model.PostgresConfig) int { return v.MaxIdleConns }),
		intValue(config.Postgres, func(v *model.PostgresConfig) int { return v.ConnMaxLifetimeMinutes }),
	)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return shutdown, err
	}
	store.InitPostgres(db)
	api.SetVerifySharedToken(stringValue(config.Auth, func(v *model.AuthConfig) string { return v.SharedToken }))
	return func(shutdownCtx context.Context) error {
		closeErr := db.Close()
		shutdownErr := shutdown(shutdownCtx)
		if shutdownErr != nil {
			return shutdownErr
		}
		return closeErr
	}, nil
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

func intValue[T any](value *T, getter func(*T) int) int {
	if value == nil {
		return 0
	}
	return getter(value)
}
