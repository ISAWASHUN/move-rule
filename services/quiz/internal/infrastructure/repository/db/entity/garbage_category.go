package mysql

import "github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/infrastructure/repository/db/mysql"

type GarbageCategory struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	mysql.Timestamp
}
