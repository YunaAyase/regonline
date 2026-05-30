package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Classes  ClassesConfig  `mapstructure:"classes"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver      string `mapstructure:"driver"`
	Path        string `mapstructure:"path"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
}

type StorageConfig struct {
	PhotoDir string `mapstructure:"photo_dir"`
}

type ClassesConfig struct {
	SeedEnabled bool            `mapstructure:"seed_enabled"`
	Presets     []ClassPreset   `mapstructure:"presets"`
}

type ClassPreset struct {
	Name        string `mapstructure:"name"`
	MaxStudents int    `mapstructure:"max_students"`
	MinAge      int    `mapstructure:"min_age"`
	MaxAge      int    `mapstructure:"max_age"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.path", "./data/regonline.db")
	v.SetDefault("database.auto_migrate", true)
	v.SetDefault("storage.photo_dir", "./photos")
	v.SetDefault("classes.seed_enabled", true)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if !filepath.IsAbs(cfg.Storage.PhotoDir) {
		if abs, err := filepath.Abs(cfg.Storage.PhotoDir); err == nil {
			cfg.Storage.PhotoDir = abs
		}
	}

	if !filepath.IsAbs(cfg.Database.Path) {
		if abs, err := filepath.Abs(cfg.Database.Path); err == nil {
			cfg.Database.Path = abs
		}
	}

	return &cfg, nil
}

func EnsureDirs(cfg *Config) error {
	dirs := []string{
		filepath.Dir(cfg.Database.Path),
		cfg.Storage.PhotoDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
