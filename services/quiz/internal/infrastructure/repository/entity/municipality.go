package mysql

type Municipality struct {
	ID   int    `gorm:"primaryKey"`
	Code int    `gorm:"not null"`
	Name string `gorm:"not null"`
	Timestamp
}
