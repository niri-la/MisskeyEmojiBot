package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	cfg "MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/errors"
)

type S3Repository interface {
	UploadFile(key string, data []byte, contentType string) (string, error)
	DownloadFile(key string) ([]byte, error)
	DeleteFile(key string) error
	GetFileURL(key string) string
}

type s3Repository struct {
	client *s3.Client
	bucket string
	config *cfg.Config
}

func NewS3Repository(cfg *cfg.Config) (S3Repository, error) {
	if !cfg.UseS3 {
		return nil, fmt.Errorf("S3 is not enabled in configuration")
	}

	// Create AWS config with custom credentials
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, errors.Config("failed to load AWS config", err)
	}

	// Create S3 client with optional custom endpoint (for S3-compatible services)
	var client *s3.Client
	if cfg.S3Endpoint != "" {
		client = s3.NewFromConfig(awsConfig, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = cfg.S3ForcePathStyle
		})
	} else {
		client = s3.NewFromConfig(awsConfig)
	}

	return &s3Repository{
		client: client,
		bucket: cfg.S3Bucket,
		config: cfg,
	}, nil
}

func (r *s3Repository) UploadFile(key string, data []byte, contentType string) (string, error) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", errors.FileOperation("failed to upload file to S3", err)
	}

	return r.GetFileURL(key), nil
}

func (r *s3Repository) DownloadFile(key string) ([]byte, error) {
	result, err := r.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, errors.FileOperation("failed to download file from S3", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, errors.FileOperation("failed to read file data from S3", err)
	}

	return data, nil
}

func (r *s3Repository) DeleteFile(key string) error {
	_, err := r.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return errors.FileOperation("failed to delete file from S3", err)
	}

	return nil
}

func (r *s3Repository) GetFileURL(key string) string {
	if r.config.S3Endpoint != "" {
		// For custom endpoints (S3-compatible services)
		if r.config.S3ForcePathStyle {
			return fmt.Sprintf("%s/%s/%s", r.config.S3Endpoint, r.bucket, key)
		}
		return fmt.Sprintf("https://%s.%s/%s", r.bucket, r.config.S3Endpoint, key)
	}

	// For AWS S3
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", r.bucket, r.config.AWSRegion, key)
}

// Helper function to determine content type from file extension
func GetContentTypeFromExtension(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".json":
		return "application/json"
	case ".db":
		return "application/octet-stream"
	default:
		return "application/octet-stream"
	}
}
