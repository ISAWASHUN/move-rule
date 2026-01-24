package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/internal/domain"
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

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) FetchData(url string) ([]domain.GarbageItem, error) {
	var allHits []Hits

	// APIの仕様により、1000件ごとにデータを取得する
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

		response, err := c.httpClient.Do(req)
		if err != nil {
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

	result := make([]domain.GarbageItem, len(allHits))
	for i, hit := range allHits {
		result[i] = domain.GarbageItem{
			Row:                      hit.Row,
			MunicipalityCode:         hit.MunicipalityCode,
			ID:                       hit.ID,
			MunicipalityName:         hit.MunicipalityName,
			AreaName:                 hit.AreaName,
			ItemName:                 hit.ItemName,
			ItemNameKana:             hit.ItemNameKana,
			ItemNameEnglish:          hit.ItemNameEnglish,
			Category:                 hit.Category,
			Notes:                    hit.Notes,
			Remarks:                  hit.Remarks,
			BulkGarbageCollectionFee: hit.BulkGarbageCollectionFee,
		}
	}

	return result, nil
}
