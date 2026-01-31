package repository

import (
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/save-garbage-categories/internal/domain"
	"gorm.io/gorm"
)

type GarbageItemRepository struct {
	db *gorm.DB
}

func NewGarbageItemRepository(db *gorm.DB) *GarbageItemRepository {
	return &GarbageItemRepository{db: db}
}

func (r *GarbageItemRepository) Create(item domain.GarbageItem, municipalityID, garbageCategoryID int, bulkGarbageFee *int) error {
	entity := GarbageItem{
		MunicipalityID:    municipalityID,
		GarbageCategoryID: garbageCategoryID,
		AreaName:          item.AreaName,
		ItemName:          item.ItemName,
		Notes:             item.Notes,
		Remarks:           item.Remarks,
		BulkGarbageFee:    bulkGarbageFee,
	}

	return r.db.Create(&entity).Error
}

func (r *GarbageItemRepository) Truncate() error {
	return r.db.Exec("TRUNCATE TABLE garbage_items").Error
}
