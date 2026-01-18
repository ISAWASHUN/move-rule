package main

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Mode     string `env:"MODE" envDefault:"up"`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	User     string `env:"DB_USER" envDefault:"root"`
	Password string `env:"DB_PASSWORD" envDefault:"password"`
	Port     string `env:"DB_PORT" envDefault:"3306"`
	TargetDB string `env:"TARGET_DB" envDefault:"all"`
	SSLMode  string `env:"SSL_MODE" envDefault:"false"`
}

func main() {
	slog.Info("start db-migrator")

	config := &Config{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	slog.Info("mode", "mode", config.Mode)
	slog.Info("connecting to db",
		"host", config.Host,
		"user", config.User,
		"port", config.Port,
		"targetDB", config.TargetDB,
		"sslMode", config.SSLMode)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to mysql: %v", err)
	}

	operator := &operator{
		db:     db,
		config: config,
	}

	// ターゲットDBが指定されている場合はそのDBに対して実行する
	if config.TargetDB != "" && config.TargetDB != "all" {
		if err := operator.execute(config.TargetDB); err != nil {
			log.Fatalf("failed to execute migrations: %v", err)
		}
		return
	}

	if config.Mode == "down" {
		log.Fatalf("rollback is not allowed for all databases")
	}

	// ターゲットDBが指定されていない場合はすべてのDBに対して実行
	dirs, err := os.ReadDir("db")
	if err != nil {
		log.Fatalf("failed to read db directory: %v", err)
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			dbName := dir.Name()
			if err := operator.execute(dbName); err != nil {
				log.Fatalf("failed to execute migrations: %v", err)
			}
		}
	}
	slog.Info("connected to default db")
}

type operator struct {
	db     *gorm.DB
	config *Config
}

func (o *operator) execute(dbName string) error {
	var exists bool
	err := o.db.Raw("SELECT EXISTS(SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?)", dbName).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %v", err)
	}

	if !exists {
		if err := o.db.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)).Error; err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
	}

	cfg := o.config
	var dsn string
	if cfg.SSLMode == "true" {
		dsn = fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?tls=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName)
	} else {
		dsn = fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?tls=false", cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName)
	}
	m, err := migrate.New(
		fmt.Sprintf("file://db/%s/migrations", dbName),
		dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	switch cfg.Mode {
	case "up":
		slog.Info("start to apply migrations", "dbName", dbName)
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to apply migrations: %v", err)
		}
		slog.Info("finish to apply migrations")
	case "down":
		slog.Info("start to rollback migrations", "dbName", dbName)
		if err := m.Steps(-1); err != nil {
			return fmt.Errorf("failed to rollback migrations: %v", err)
		}
		slog.Info("finish to rollback migrations", "dbName", dbName)
	default:
		return fmt.Errorf("unknown mode: %s", cfg.Mode)
	}
	return nil
}
