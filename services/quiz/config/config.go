package config

import (
	"os"

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
