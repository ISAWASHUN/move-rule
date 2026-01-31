package usecase

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/save-garbage-categories/internal/domain"
)

type SaveGarbageCategoriesUseCase struct {
	fileReader          domain.FileReader
	municipalityRepo    domain.MunicipalityRepository
	garbageCategoryRepo domain.GarbageCategoryRepository
	garbageItemRepo     domain.GarbageItemRepository
}

func NewSaveGarbageCategoriesUseCase(
	fileReader domain.FileReader,
	municipalityRepo domain.MunicipalityRepository,
	garbageCategoryRepo domain.GarbageCategoryRepository,
	garbageItemRepo domain.GarbageItemRepository,
) *SaveGarbageCategoriesUseCase {
	return &SaveGarbageCategoriesUseCase{
		fileReader:          fileReader,
		municipalityRepo:    municipalityRepo,
		garbageCategoryRepo: garbageCategoryRepo,
		garbageItemRepo:     garbageItemRepo,
	}
}

func (u *SaveGarbageCategoriesUseCase) Execute(filePath string) error {
	items, err := u.fileReader.Read(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	log.Printf("read %d items", len(items))

	// テーブルをTRUNCATE（外部キー制約の順序に注意）
	log.Println("truncating tables...")
	if err := u.garbageItemRepo.Truncate(); err != nil {
		return fmt.Errorf("failed to truncate garbage_items: %w", err)
	}
	if err := u.garbageCategoryRepo.Truncate(); err != nil {
		return fmt.Errorf("failed to truncate garbage_categories: %w", err)
	}
	if err := u.municipalityRepo.Truncate(); err != nil {
		return fmt.Errorf("failed to truncate municipalities: %w", err)
	}

	log.Println("saving data...")
	for i, item := range items {
		if err := u.saveItem(item); err != nil {
			return fmt.Errorf("failed to save item at row %d: %w", i+1, err)
		}

		if (i+1)%1000 == 0 {
			log.Printf("progress: %d/%d", i+1, len(items))
		}
	}

	log.Printf("saved %d items successfully", len(items))
	return nil
}

func (u *SaveGarbageCategoriesUseCase) saveItem(item domain.GarbageItem) error {
	municipalityID, err := u.municipalityRepo.GetOrCreate(item.MunicipalityCode, item.MunicipalityName)
	if err != nil {
		return fmt.Errorf("failed to get or create municipality: %w", err)
	}

	category := item.Category
	if category == "" {
		category = "未分類"
	}
	garbageCategoryID, err := u.garbageCategoryRepo.GetOrCreate(category)
	if err != nil {
		return fmt.Errorf("failed to get or create garbage category: %w", err)
	}

	var bulkGarbageFee *int
	if item.BulkGarbageCollectionFee != "" {
		fee, err := strconv.Atoi(item.BulkGarbageCollectionFee)
		if err == nil {
			bulkGarbageFee = &fee
		}
	}

	if err := u.garbageItemRepo.Create(item, municipalityID, garbageCategoryID, bulkGarbageFee); err != nil {
		return fmt.Errorf("failed to create garbage item: %w", err)
	}

	return nil
}
