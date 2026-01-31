package mysql

type GarbageItem struct {
	ID                int    `gorm:"primaryKey"`
	MunicipalityID    int    `gorm:"not null"`
	GarbageCategoryID int    `gorm:"not null"`
	AreaName          string `gorm:"size:255"`
	ItemName          string `gorm:"size:255;not null"`
	Notes             string `gorm:"type:text"`
	Remarks           string `gorm:"type:text"`
	BulkGarbageFee    int    `gorm:"type:int"`
	Timestamp
}
