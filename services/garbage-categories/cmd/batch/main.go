package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	baseUrl      = "https://service.api.metro.tokyo.lg.jp"
	itabashiUrl  = baseUrl + "/api/t131199d3000000001-10af70080e2503877feb2bf2c9a42171-0/json"
)

type APIResponse struct {
	Total    int      `json:"total"`
	Subtotal int      `json:"subtotal"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Metadata Metadata `json:"metadata"`
	Hits     []Hits   `json:"hits"`
}

type Metadata struct {
	APIID        string `json:"apiId"`
	Title        string `json:"title"`
	DatasetID    string `json:"datasetId"`
	DatasetTitle string `json:"datasetTitle"`
	DatasetDesc  string `json:"datasetDesc"`
	DataTitle    string `json:"dataTitle"`
	DataDesc     string `json:"dataDesc"`
	SheetName    string `json:"sheetname"`
	Version      string `json:"version"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
}

type Hits struct {
	Row                      int    `json:"row"`
	MunicipalityCode         int    `json:"全国地方公共団体コード"`
	ID                       string `json:"ID"`
	MunicipalityName         string `json:"地方公共団体名"`
	AreaName                 string `json:"地区名"`
	ItemName                 string `json:"ゴミの品目"`
	ItemNameKana             string `json:"ゴミの品目_カナ"`
	ItemNameEnglish          string `json:"ゴミの品目_英字"`
	Category                 string `json:"分別区分"`
	Notes                    string `json:"注意点"`
	Remarks                  string `json:"備考"`
	BulkGarbageCollectionFee string `json:"粗大ごみ回収料金"`
}

type Municipality struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Code      int       `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type GarbageCategory struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:255;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type GarbageItem struct {
	ID              int       `gorm:"primaryKey;autoIncrement"`
	MunicipalityID  int       `gorm:"not null"`
	GarbageCategoryID int       `gorm:"not null"`
	AreaName        string    `gorm:"size:255"`
	ItemName        string    `gorm:"size:255;not null"`
	Notes           string    `gorm:"type:text"`
	Remarks         string    `gorm:"type:text"`
	BulkGarbageFee  *int      `gorm:"type:int"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := truncateTables(db); err != nil {
		log.Fatalf("failed to truncate tables: %v", err)
	}

	urls := []string{itabashiUrl}
	var allHits []Hits

	for _, url := range urls {
		log.Printf("fetch data from %s", url)
		hits, err := fetchData(url)
		if err != nil {
			log.Fatalf("failed to fetch data from %s: %v", url, err)
		}
		allHits = append(allHits, hits...)
	}

	log.Println("save data...")
	if err := saveData(db, allHits); err != nil {
		log.Fatalf("failed to save data: %v", err)
	}

	log.Println("data saved successfully")
}

func connectDB() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "garbage_category_rule_quiz")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func truncateTables(db *gorm.DB) error {
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	tables := []string{"garbage_items", "garbage_categories", "municipalities"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)).Error; err != nil {
			return err
		}
	}

	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return err
	}

	return nil
}

func fetchData(url string) ([]Hits, error) {
	var allHits []Hits
	limit := 1000
	offset := 0

	for {
		paginatedURL := fmt.Sprintf("%s?limit=%d&offset=%d", url, limit, offset)

		reqBody := bytes.NewBuffer([]byte(`{}`))

		req, err := http.NewRequest("POST", paginatedURL, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			// リトライなし、一時的なネットワークエラーでも即座に失敗
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}

		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}

		respBody, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var data APIResponse
		if err := json.Unmarshal(respBody, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allHits = append(allHits, data.Hits...)

		if len(allHits) >= data.Total {
			break
		}

		offset += limit
		log.Printf("progress: %d/%d", len(allHits), data.Total)
	}

	return allHits, nil
}

func saveData(db *gorm.DB, hits []Hits) error {
	municipalityCache := make(map[int]int)
	garbageCategoryCache := make(map[string]int)

	for _, hit := range hits {
		municipalityID, err := getOrCreateMunicipality(db, hit.MunicipalityCode, hit.MunicipalityName, municipalityCache)
		if err != nil {
			return fmt.Errorf("failed to save municipality: %w", err)
		}

		garbageCategoryID, err := getOrCreateGarbageCategory(db, hit.Category, garbageCategoryCache)
		if err != nil {
			return fmt.Errorf("failed to save garbage category: %w", err)
		}

		var bulkGarbageFee *int
		if hit.BulkGarbageCollectionFee != "" {
			fee, err := strconv.Atoi(hit.BulkGarbageCollectionFee)
			if err == nil {
				bulkGarbageFee = &fee
			}
		}

		garbageItem := GarbageItem{
			MunicipalityID:   municipalityID,
			GarbageCategoryID: garbageCategoryID,
			AreaName:        hit.AreaName,
			ItemName:        hit.ItemName,
			Notes:           hit.Notes,
			Remarks:         hit.Remarks,
			BulkGarbageFee:  bulkGarbageFee,
		}

		if err := db.Create(&garbageItem).Error; err != nil {
			return fmt.Errorf("failed to create garbage item: %w", err)
		}
	}

	return nil
}

func getOrCreateMunicipality(db *gorm.DB, code int, name string, cache map[int]int) (int, error) {
	if id, exists := cache[code]; exists {
		return id, nil
	}

	municipality := Municipality{
		Code: code,
		Name: name,
	}

	if err := db.Create(&municipality).Error; err != nil {
		return 0, err
	}

	cache[code] = municipality.ID
	return municipality.ID, nil
}

func getOrCreateGarbageCategory(db *gorm.DB, name string, cache map[string]int) (int, error) {
	if name == "" {
		name = "未分類"
	}

	if id, exists := cache[name]; exists {
		return id, nil
	}

	garbageCategory := GarbageCategory{
		Name: name,
	}

	if err := db.Create(&garbageCategory).Error; err != nil {
		return 0, err
	}

	cache[name] = garbageCategory.ID
	return garbageCategory.ID, nil
}
