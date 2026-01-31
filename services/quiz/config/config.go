package config

import (
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
)

type Server struct {
	Host string `toml:"host" validate:"required,hostname|ip"`
	Port int    `toml:"port" validate:"required,min=1,max=65535"`
}

type App struct {
	LogLevel string `toml:"log_level" validate:"required,oneof=debug info warn error"`
}

type MySQL struct {
	Host     string `toml:"host" validate:"required"`
	Port     string `toml:"port" validate:"required"`
	User     string `toml:"user" validate:"required"`
	Password string `toml:"password" validate:"required"`
	DBName   string `toml:"db_name" validate:"required"`
	SSLMode  string `toml:"ssl_mode" validate:"required,oneof=disable require"`
	UseMock  bool   `toml:"use_mock"`
}

type Config struct {
	Server Server `toml:"server" validate:"required"`
	App    App    `toml:"app" validate:"required"`
	MySQL  MySQL  `toml:"mysql" validate:"required"`
}

var validate = validator.New()

// Load はTOMLファイルから設定を読み込みます（ローカル開発用）
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := toml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadFromEnv は環境変数から設定を読み込みます（Lambda用）
func LoadFromEnv() *Config {
	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))

	return &Config{
		Server: Server{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: port,
		},
		App: App{
			LogLevel: getEnv("APP_LOG_LEVEL", "info"),
		},
		MySQL: MySQL{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "garbage_category_rule_quiz"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			UseMock:  getEnv("DB_USE_MOCK", "false") == "true",
		},
	}
}

// getEnv は環境変数を取得します。存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
