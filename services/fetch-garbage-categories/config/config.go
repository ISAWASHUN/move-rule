package config

import "os"

type Config struct {
	S3Bucket string
	S3Prefix string
	AWSRegion string
	Env string
}

func Load() *Config {
	return &Config{
		S3Bucket:  getEnv("S3_BUCKET", ""),
		S3Prefix:  getEnv("S3_PREFIX", ""),
		AWSRegion: getEnv("AWS_REGION", "ap-northeast-1"),
		Env:       getEnv("ENV", "local"),
	}
}

func (c *Config) IsS3Enabled() bool {
	return c.S3Bucket != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
