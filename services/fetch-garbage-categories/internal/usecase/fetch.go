package usecase

import (
	"log"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/internal/domain"
)

type FetchGarbageCategoriesUseCase struct {
	apiClient domain.APIClient
	storage   domain.Storage
}

func NewFetchGarbageCategoriesUseCase(
	apiClient domain.APIClient,
	storage domain.Storage,
) *FetchGarbageCategoriesUseCase {
	return &FetchGarbageCategoriesUseCase{
		apiClient: apiClient,
		storage:   storage,
	}
}

func (u *FetchGarbageCategoriesUseCase) Execute(urls []string) error {
	var allItems []domain.GarbageItem

	for _, url := range urls {
		log.Printf("fetching data from %s", url)
		items, err := u.apiClient.FetchData(url)
		if err != nil {
			return err
		}
		allItems = append(allItems, items...)
		log.Printf("fetched %d items from %s", len(items), url)
	}

	log.Printf("total items fetched: %d", len(allItems))

	if err := u.storage.Save(allItems); err != nil {
		return err
	}

	log.Printf("data saved successfully")
	return nil
}
