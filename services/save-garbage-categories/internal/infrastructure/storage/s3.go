package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/save-garbage-categories/internal/domain"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3FileReader はS3上のJSONファイルを読み込むためのリーダーです。
type S3FileReader struct {
	client *s3.Client
	bucket string
}

// NewS3FileReader はS3FileReaderを初期化します。
// bucket はS3_BUCKET、region はAWS_REGIONの想定です。
func NewS3FileReader(ctx context.Context, bucket, region string) (*S3FileReader, error) {
	if region == "" {
		region = "ap-northeast-1"
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3FileReader{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}, nil
}

// Read はS3上のJSONファイルを読み込みます。
// path はS3キー、または s3://bucket/key 形式を受け付けます。
func (r *S3FileReader) Read(path string) ([]domain.GarbageItem, error) {
	bucket := r.bucket
	key := strings.TrimPrefix(path, "/")
	if strings.HasPrefix(path, "s3://") {
		parsedBucket, parsedKey, err := parseS3Path(path)
		if err != nil {
			return nil, err
		}
		bucket = parsedBucket
		key = parsedKey
	}

	if bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}

	resp, err := r.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3 (bucket=%s, key=%s): %w", bucket, key, err)
	}
	defer resp.Body.Close()

	var items []domain.GarbageItem
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&items); err != nil {
		return nil, fmt.Errorf("failed to decode JSON from S3: %w", err)
	}

	return items, nil
}

func parseS3Path(path string) (string, string, error) {
	trimmed := strings.TrimPrefix(path, "s3://")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid s3 path: %s", path)
	}
	return parts[0], parts[1], nil
}
