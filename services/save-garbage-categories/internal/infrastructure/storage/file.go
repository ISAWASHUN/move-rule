package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/save-garbage-categories/internal/domain"
)

type JSONFileReader struct{}

func NewJSONFileReader() *JSONFileReader {
	return &JSONFileReader{}
}

func (r *JSONFileReader) Read(path string) ([]domain.GarbageItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var items []domain.GarbageItem
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return items, nil
}
