package mysql

import "time"

// Timestamp は作成日時と更新日時を管理する構造体です
type Timestamp struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type GarbageCategory struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	Timestamp
}
