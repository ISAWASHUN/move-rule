package mysql

import (
	"fmt"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	logLevelMap = map[string]logger.LogLevel{
		"debug": logger.Info,
		"info":  logger.Silent,
		"warn":  logger.Warn,
		"error": logger.Error,
	}
)

func ConnectDB(cfg config.MySQL, logLevel string) (*gorm.DB, error) {
	logger := logger.Default.LogMode(logLevelMap[logLevel])

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

type Timestamp struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
