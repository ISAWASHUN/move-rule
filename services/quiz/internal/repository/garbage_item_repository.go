package repository

import (
	"context"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/domain"
	models "github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/repository/db/models"
	"gorm.io/gorm"
)

type GarbageItemRepository interface {
	domain.GarbageItemRepository
}

type garbageItemRepository struct {
	db *gorm.DB
}

func NewGarbageItemRepository(db *gorm.DB) GarbageItemRepository {
	return &garbageItemRepository{db: db}
}

func (r *garbageItemRepository) GetByMunicipalityID(ctx context.Context, municipalityID int) ([]domain.GarbageItem, error) {
	var items []models.GarbageItem
	err := r.db.WithContext(ctx).
		Where("municipality_id = ?", municipalityID).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.GarbageItem, len(items))
	for i, item := range items {
		result[i] = toDomainGarbageItem(item)
	}
	return result, nil
}

func (r *garbageItemRepository) GetByID(ctx context.Context, id int) (*domain.GarbageItem, error) {
	var item models.GarbageItem
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, err
	}

	result := toDomainGarbageItem(item)
	return &result, nil
}

func (r *garbageItemRepository) GetByIDWithCategory(ctx context.Context, id int) (*domain.GarbageItem, *domain.GarbageCategory, error) {
	var item models.GarbageItem
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, nil, err
	}

	var category models.GarbageCategory
	err = r.db.WithContext(ctx).First(&category, item.GarbageCategoryID).Error
	if err != nil {
		return nil, nil, err
	}

	domainItem := toDomainGarbageItem(item)
	domainCategory := toDomainGarbageCategory(category)
	return &domainItem, &domainCategory, nil
}

func toDomainGarbageItem(m models.GarbageItem) domain.GarbageItem {
	return domain.GarbageItem{
		ID:                domain.GarbageItemID(m.ID),
		MunicipalityID:    domain.MunicipalityID(m.MunicipalityID),
		GarbageCategoryID: domain.GarbageCategoryID(m.GarbageCategoryID),
		AreaName:          m.AreaName,
		ItemName:          m.ItemName,
		Notes:             m.Notes,
		Remarks:           m.Remarks,
		BulkGarbageFee:    m.BulkGarbageFee,
	}
}

func toDomainGarbageCategory(m models.GarbageCategory) domain.GarbageCategory {
	return domain.GarbageCategory{
		ID:   domain.GarbageCategoryID(m.ID),
		Name: m.Name,
	}
}
