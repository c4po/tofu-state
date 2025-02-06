package storage

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

func NewStorageBackend(cfg StorageConfig) StorageBackend {
	switch cfg.Type {
	case "s3":
		return NewS3Storage(cfg)
	case "gcs":
		return NewGCSStorage(cfg)
	default:
		return NewLocalStorage(cfg)
	}
}
