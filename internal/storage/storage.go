package storage

import "go.uber.org/zap"

type StorageBackend interface {
	GetState(workspace string) ([]byte, error)
	PutState(workspace string, data []byte) error
	ListModules() ([]string, error)
	UploadModule(name string, data []byte) error
}

type StorageConfig struct {
	Type       string
	BucketName string
	Region     string // For S3
	LocalPath  string // For local storage
}

func NewStorageBackend(cfg StorageConfig, logger *zap.Logger) StorageBackend {
	switch cfg.Type {
	case "s3":
		return NewS3Storage(cfg, logger)
	case "gcs":
		return NewGCSStorage(cfg, logger)
	default:
		return NewLocalStorage(cfg, logger)
	}
}
