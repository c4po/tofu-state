package storage

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

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

func NewStorageBackend(logger *zap.Logger) StorageBackend {
	cfg := StorageConfig{
		Type:       viper.GetString("storage.type"),
		BucketName: viper.GetString("storage.bucket_name"),
		Region:     viper.GetString("storage.region"),
		LocalPath:  viper.GetString("storage.local_path"),
	}
	switch cfg.Type {
	case "s3":
		return NewS3Storage(cfg, logger)
	case "gcs":
		return NewGCSStorage(cfg, logger)
	default:
		return NewLocalStorage(cfg, logger)
	}
}
