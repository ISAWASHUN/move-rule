package repository

import "gorm.io/gorm"

type MunicipalityRepository struct {
	db    *gorm.DB
	cache map[int]int
}

func NewMunicipalityRepository(db *gorm.DB) *MunicipalityRepository {
	return &MunicipalityRepository{
		db:    db,
		cache: make(map[int]int),
	}
}

func (r *MunicipalityRepository) GetOrCreate(code int, name string) (int, error) {
	if id, exists := r.cache[code]; exists {
		return id, nil
	}

	municipality := Municipality{
		Code: code,
		Name: name,
	}

	if err := r.db.Create(&municipality).Error; err != nil {
		return 0, err
	}

	r.cache[code] = municipality.ID
	return municipality.ID, nil
}

func (r *MunicipalityRepository) Truncate() error {
	return r.db.Exec("TRUNCATE TABLE municipalities").Error
}
