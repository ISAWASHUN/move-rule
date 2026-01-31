package repository

import "time"

// Timestamp は作成日時と更新日時を管理する構造体です
type Timestamp struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Municipality struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Code int    `gorm:"uniqueIndex;not null"`
	Name string `gorm:"size:255;not null"`
	Timestamp
}

func (Municipality) TableName() string {
	return "municipalities"
}

type GarbageCategory struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:255;uniqueIndex;not null"`
	Timestamp
}

func (GarbageCategory) TableName() string {
	return "garbage_categories"
}

type GarbageItem struct {
	ID                int    `gorm:"primaryKey;autoIncrement"`
	MunicipalityID    int    `gorm:"not null"`
	GarbageCategoryID int    `gorm:"not null"`
	AreaName          string `gorm:"size:255"`
	ItemName          string `gorm:"size:255;not null"`
	Notes             string `gorm:"type:text"`
	Remarks           string `gorm:"type:text"`
	BulkGarbageFee    *int   `gorm:"type:int"`
	Timestamp
}

func (GarbageItem) TableName() string {
	return "garbage_items"
}
