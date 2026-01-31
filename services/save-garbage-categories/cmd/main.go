package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/save-garbage-categories/internal/domain"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/save-garbage-categories/internal/infrastructure/repository"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/save-garbage-categories/internal/infrastructure/storage"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/save-garbage-categories/internal/usecase"
)

var (
	logLevelMap = map[string]logger.LogLevel{
		"debug": logger.Info,
		"info":  logger.Silent,
		"warn":  logger.Warn,
		"error": logger.Error,
	}
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

func handler(ctx context.Context) (Response, error) {
	if err := run(ctx); err != nil {
		return Response{
			Message: err.Error(),
			Status:  "error",
		}, err
	}

	return Response{
		Message: "data saved to database successfully",
		Status:  "success",
	}, nil
}

func run(ctx context.Context) error {
	filePath := getEnv("INPUT_FILE", "../fetch-garbage-categories/internal/infrastructure/storage/file/latest.json")
	s3Bucket := getEnv("S3_BUCKET", "")
	s3Key := getEnv("S3_KEY", "")
	s3Prefix := getEnv("S3_PREFIX", "")
	awsRegion := getEnv("AWS_REGION", "ap-northeast-1")

	db, err := connectDBFromEnv()
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 外部キー制約を一時的に無効化
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}
	defer func() {
		if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
			log.Printf("failed to enable foreign key checks: %v", err)
		}
	}()

	fileReader, filePath, err := buildFileReader(ctx, filePath, s3Bucket, s3Key, s3Prefix, awsRegion)
	if err != nil {
		return fmt.Errorf("failed to build file reader: %w", err)
	}
	municipalityRepo := repository.NewMunicipalityRepository(db)
	garbageCategoryRepo := repository.NewGarbageCategoryRepository(db)
	garbageItemRepo := repository.NewGarbageItemRepository(db)

	saveUseCase := usecase.NewSaveGarbageCategoriesUseCase(
		fileReader,
		municipalityRepo,
		garbageCategoryRepo,
		garbageItemRepo,
	)

	if err := saveUseCase.Execute(filePath); err != nil {
		return fmt.Errorf("failed to execute usecase: %w", err)
	}

	log.Println("data saved to database successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func connectDBFromEnv() (*gorm.DB, error) {
	cfg := struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "garbage_category_rule_quiz"),
	}

	logLevel := getEnv("LOG_LEVEL", "info")
	level := logLevelMap[logLevel]

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func buildFileReader(
	ctx context.Context,
	filePath string,
	s3Bucket string,
	s3Key string,
	s3Prefix string,
	awsRegion string,
) (domain.FileReader, string, error) {
	useS3 := s3Bucket != "" || s3Key != "" || strings.HasPrefix(filePath, "s3://")
	if !useS3 {
		return storage.NewJSONFileReader(), filePath, nil
	}

	reader, err := storage.NewS3FileReader(ctx, s3Bucket, awsRegion)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create s3 reader: %w", err)
	}

	if s3Key == "" && !strings.HasPrefix(filePath, "s3://") {
		s3Key = defaultS3Key(s3Prefix)
	}

	if s3Key != "" && !strings.HasPrefix(filePath, "s3://") {
		filePath = s3Key
	}

	return reader, filePath, nil
}

func defaultS3Key(prefix string) string {
	prefix = strings.TrimPrefix(prefix, "/")
	if prefix == "" {
		return "data/latest.json"
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	return prefix + "latest.json"
}
