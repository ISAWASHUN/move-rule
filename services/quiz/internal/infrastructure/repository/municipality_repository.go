package repository

import (
	"context"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/domain"
	entity "github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/infrastructure/repository/entity"
	"gorm.io/gorm"
)

type MunicipalityRepository interface {
	domain.MunicipalityRepository
}

type municipalityRepository struct {
	db *gorm.DB
}

func NewMunicipalityRepository(db *gorm.DB) MunicipalityRepository {
	return &municipalityRepository{db: db}
}

func (r *municipalityRepository) GetByName(ctx context.Context, name string) (*domain.Municipality, error) {
	var municipality entity.Municipality
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&municipality).Error
	if err != nil {
		return nil, err
	}

	return &domain.Municipality{
		ID:   domain.MunicipalityID(municipality.ID),
		Code: municipality.Code,
		Name: municipality.Name,
	}, nil
}

func (r *municipalityRepository) GetByID(ctx context.Context, id int) (*domain.Municipality, error) {
	var municipality entity.Municipality
	err := r.db.WithContext(ctx).First(&municipality, id).Error
	if err != nil {
		return nil, err
	}

	return &domain.Municipality{
		ID:   domain.MunicipalityID(municipality.ID),
		Code: municipality.Code,
		Name: municipality.Name,
	}, nil
}

func (r *municipalityRepository) GetByCode(ctx context.Context, code int) (*domain.Municipality, error) {
	var municipality entity.Municipality
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&municipality).Error
	if err != nil {
		return nil, err
	}

	return &domain.Municipality{
		ID:   domain.MunicipalityID(municipality.ID),
		Code: municipality.Code,
		Name: municipality.Name,
	}, nil
}
