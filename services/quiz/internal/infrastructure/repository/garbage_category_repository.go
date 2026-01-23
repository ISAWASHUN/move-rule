package repository

import (
	"context"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/domain"
	models "github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/infrastructure/repository/db/entity"
	"gorm.io/gorm"
)

type GarbageCategoryRepository interface {
	domain.GarbageCategoryRepository
}

type garbageCategoryRepository struct {
	db *gorm.DB
}

func NewGarbageCategoryRepository(db *gorm.DB) GarbageCategoryRepository {
	return &garbageCategoryRepository{db: db}
}

func (r *garbageCategoryRepository) GetAll(ctx context.Context) ([]domain.GarbageCategory, error) {
	var categories []models.GarbageCategory
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.GarbageCategory, len(categories))
	for i, cat := range categories {
		result[i] = domain.GarbageCategory{
			ID:   domain.GarbageCategoryID(cat.ID),
			Name: cat.Name,
		}
	}
	return result, nil
}

func (r *garbageCategoryRepository) GetByName(ctx context.Context, name string) (*domain.GarbageCategory, error) {
	var category models.GarbageCategory
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}

	return &domain.GarbageCategory{
		ID:   domain.GarbageCategoryID(category.ID),
		Name: category.Name,
	}, nil
}

func (r *garbageCategoryRepository) GetByID(ctx context.Context, id int) (*domain.GarbageCategory, error) {
	var category models.GarbageCategory
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}

	return &domain.GarbageCategory{
		ID:   domain.GarbageCategoryID(category.ID),
		Name: category.Name,
	}, nil
}
