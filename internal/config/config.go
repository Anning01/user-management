package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey       string
	ExpirationHours int
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 读取环境变量
	viper.AutomaticEnv()

	// 默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.readTimeout", 5)
	viper.SetDefault("server.writeTimeout", 10)
	viper.SetDefault("server.maxHeaderBytes", 1<<20)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
