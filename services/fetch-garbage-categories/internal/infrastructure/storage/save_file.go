package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/internal/domain"
)

type JSONStorage struct {
	outputDir string
}

func NewJSONStorage(outputDir string) *JSONStorage {
	return &JSONStorage{
		outputDir: outputDir,
	}
}

func (s *JSONStorage) Save(items []domain.GarbageItem) error {
	if err := os.MkdirAll(s.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := fmt.Sprintf("garbage_categories_%s.json", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(s.outputDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(items); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	latestPath := filepath.Join(s.outputDir, "latest.json")
	latestFile, err := os.Create(latestPath)
	if err != nil {
		return fmt.Errorf("failed to create latest file: %w", err)
	}
	defer latestFile.Close()

	latestEncoder := json.NewEncoder(latestFile)
	latestEncoder.SetIndent("", "  ")
	if err := latestEncoder.Encode(items); err != nil {
		return fmt.Errorf("failed to encode latest JSON: %w", err)
	}

	return nil
}
