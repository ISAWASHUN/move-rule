package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/fetch-garbage-categories/internal/domain"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage は domain.Storage を S3 で実装したものです。
// 環境変数で指定されたバケット / プレフィックスに JSON を保存します。
type S3Storage struct {
	client *s3.Client
	bucket string
	prefix string
}

// NewS3Storage は S3Storage を初期化します。
// region は AWS_REGION、bucket は S3_BUCKET、prefix は S3_PREFIX などから呼び出し元で渡してください。
func NewS3Storage(ctx context.Context, bucket, region, prefix string) (*S3Storage, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// prefix が空でなければ末尾にスラッシュを付けておく
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix = prefix + "/"
	}

	return &S3Storage{
		client: client,
		bucket: bucket,
		prefix: prefix,
	}, nil
}

// Save は items を JSON 形式で S3 に保存します。
// タイムスタンプ付きのキーと latest.json の 2 つを書き込みます。
func (s *S3Storage) Save(items []domain.GarbageItem) error {
	ctx := context.Background()

	// JSON エンコード
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(items); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	body := bytes.NewReader(buf.Bytes())

	timestamp := time.Now().Format("20060102_150405")
	objectKey := fmt.Sprintf("%sgarbage_categories_%s.json", s.prefix, timestamp)
	latestKey := fmt.Sprintf("%slatest.json", s.prefix)

	// タイムスタンプ付きオブジェクトの保存
	if err := s.putObject(ctx, objectKey, body); err != nil {
		return err
	}

	// latest.json 用にリーダーを作り直す
	body = bytes.NewReader(buf.Bytes())
	if err := s.putObject(ctx, latestKey, body); err != nil {
		return err
	}

	return nil
}

func (s *S3Storage) putObject(ctx context.Context, key string, body *bytes.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   body,
		ACL:    types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return fmt.Errorf("failed to put object to S3 (key=%s): %w", key, err)
	}
	return nil
}
