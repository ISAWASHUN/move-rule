package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/config"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/internal/domain"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/internal/infrastructure/api"
	s3storage "github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/internal/infrastructure/storage"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/internal/usecase"
)

const (
	baseUrl     = "https://service.api.metro.tokyo.lg.jp"
	itabashiUrl = baseUrl + "/api/t131199d3000000001-10af70080e2503877feb2bf2c9a42171-0/json"
)

// Response はLambdaのレスポンス
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	// Lambda環境かどうかを判定
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda環境の場合はハンドラーを起動
		lambda.Start(handler)
	} else {
		// ローカル環境の場合は直接実行
		if err := run(context.Background()); err != nil {
			log.Fatalf("failed to run: %v", err)
		}
	}
}

// handler はLambdaハンドラー関数
func handler(ctx context.Context) (Response, error) {
	if err := run(ctx); err != nil {
		return Response{
			Message: err.Error(),
			Status:  "error",
		}, err
	}

	return Response{
		Message: "data fetched and saved successfully",
		Status:  "success",
	}, nil
}

// run はメインのビジネスロジック
func run(ctx context.Context) error {
	cfg := config.Load()
	urls := []string{itabashiUrl}

	apiClient := api.NewClient()

	var (
		storage domain.Storage
		desc    string
	)

	if cfg.IsS3Enabled() {
		s, d, err := newS3Storage(ctx, cfg)
		if err != nil {
			return fmt.Errorf("failed to create s3 storage: %w", err)
		}
		storage, desc = s, d
	} else {
		// ローカル環境の場合はファイルに保存
		outputDir := getOutputDir()
		storage = s3storage.NewJSONStorage(outputDir)
		desc = outputDir
	}

	fetchUseCase := usecase.NewFetchGarbageCategoriesUseCase(apiClient, storage)

	if err := fetchUseCase.Execute(urls); err != nil {
		return fmt.Errorf("failed to execute usecase: %w", err)
	}

	log.Printf("data saved successfully to %s", desc)
	return nil
}

func newS3Storage(ctx context.Context, cfg *config.Config) (domain.Storage, string, error) {
	s3, err := s3storage.NewS3Storage(ctx, cfg.S3Bucket, cfg.AWSRegion, cfg.S3Prefix)
	if err != nil {
		return nil, "", err
	}

	desc := "s3://" + cfg.S3Bucket
	if cfg.S3Prefix != "" {
		desc += "/" + cfg.S3Prefix
	}

	return s3, desc, nil
}

func getOutputDir() string {
	// プロジェクトルート（go.modがある場所）を特定
	projectRoot := findProjectRoot()
	return filepath.Join(projectRoot, "internal/infrastructure/storage/file")
}

func findProjectRoot() string {
	// カレントディレクトリからgo.modを探す
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}

	currentDir := wd
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// ルートディレクトリに到達した場合はカレントディレクトリを返す
			return wd
		}
		currentDir = parentDir
	}
}
