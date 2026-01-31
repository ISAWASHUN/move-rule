package repository

import "gorm.io/gorm"

type GarbageCategoryRepository struct {
	db    *gorm.DB
	cache map[string]int
}

func NewGarbageCategoryRepository(db *gorm.DB) *GarbageCategoryRepository {
	return &GarbageCategoryRepository{
		db:    db,
		cache: make(map[string]int),
	}
}

func (r *GarbageCategoryRepository) GetOrCreate(name string) (int, error) {
	if id, exists := r.cache[name]; exists {
		return id, nil
	}

	category := GarbageCategory{
		Name: name,
	}

	if err := r.db.Create(&category).Error; err != nil {
		return 0, err
	}

	r.cache[name] = category.ID
	return category.ID, nil
}

func (r *GarbageCategoryRepository) Truncate() error {
	return r.db.Exec("TRUNCATE TABLE garbage_categories").Error
}
