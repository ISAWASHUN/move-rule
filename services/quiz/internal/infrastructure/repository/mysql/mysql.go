package mysql

import (
	"fmt"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var logLevelMap = map[string]logger.LogLevel{
	"debug": logger.Info,
	"info":  logger.Silent,
	"warn":  logger.Warn,
	"error": logger.Error,
}

func NewMySQL(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logLevelMap[cfg.App.LogLevel]

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
