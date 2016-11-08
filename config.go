package pubsub

import (
	"context"

	"github.com/BurntSushi/toml"
)

type private struct{}

var cfgKey private

type Config struct {
	RedisURI string     `toml:"redis_uri"`
	GCP      *GCPConfig `toml:"gcp"`
}

type GCPConfig struct {
	ProjectID          string `toml:"project_id"`
	ServiceAccountFile string `toml:"service_account_file"`
}

func DefaultConfig() *Config {
	return &Config{
		RedisURI: "redis://127.0.0.1:6379",
		GCP: &GCPConfig{
			ProjectID:          "",
			ServiceAccountFile: "",
		},
	}
}

func NewConfig(configPath string) (*Config, error) {
	cfg := DefaultConfig()
	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewContext(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, cfgKey, cfg)
}

func ConfigFromContext(ctx context.Context) *Config {
	if cfg, ok := ctx.Value(cfgKey).(*Config); ok && cfg != nil {
		return cfg
	}
	return DefaultConfig()
}
