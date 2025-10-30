package config

import (
	"time"

	"github.com/spf13/viper"
)

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

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type JWTConfig struct {
	SecretKey       string
	ExpirationHours int
}

func Load() (*Config, error) {
	// 1. 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.readTimeout", 5)
	viper.SetDefault("server.writeTimeout", 10)
	viper.SetDefault("server.maxHeaderBytes", 1<<20)

	// 2. 先绑定环境变量（必须在读取配置文件之前）
	// 手动绑定环境变量，支持 DB_PASSWORD 这种格式
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.username", "DB_USERNAME")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.driver", "DB_DRIVER")
	viper.BindEnv("jwt.secretKey", "JWT_SECRET_KEY")
	viper.BindEnv("jwt.expirationHours", "JWT_EXPIRATION_HOURS")
	viper.BindEnv("server.port", "SERVER_PORT")

	// 3. 读取配置文件 (config.yaml)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 如果配置文件不存在，不报错，使用默认值和环境变量
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// 4. 解析为结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
