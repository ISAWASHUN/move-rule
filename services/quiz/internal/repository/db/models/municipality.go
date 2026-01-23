package mysql

import "github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/repository/db/mysql"

type Municipality struct {
	ID   int    `gorm:"primaryKey"`
	Code int    `gorm:"not null"`
	Name string `gorm:"not null"`
	mysql.Timestamp
}
