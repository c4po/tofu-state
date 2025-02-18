package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	logger *zap.Logger
}

func NewS3Storage(cfg StorageConfig, logger *zap.Logger) *S3Storage {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		logger.Error("Failed to load AWS config", zap.Error(err))
		panic(err)
	}

	return &S3Storage{
		client: s3.NewFromConfig(awsCfg),
		bucket: cfg.BucketName,
		logger: logger,
	}
}

func (s *S3Storage) GetState(workspace string) ([]byte, error) {
	s.logger.Debug("Getting state", zap.String("workspace", workspace))

	// Temporary implementation
	return nil, nil
}

func (s *S3Storage) PutState(workspace string, data []byte) error {
	s.logger.Debug("Putting state", zap.String("workspace", workspace))

	// Temporary implementation
	return nil
}

// Add missing interface methods
func (s *S3Storage) ListModules() ([]string, error) {
	s.logger.Debug("Listing modules")

	// Temporary implementation
	return []string{}, nil
}

func (s *S3Storage) UploadModule(name string, data []byte) error {
	s.logger.Debug("Uploading module", zap.String("name", name))

	// Temporary implementation
	return nil
}
