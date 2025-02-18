package storage

import "go.uber.org/zap"

type GCSStorage struct {
	logger *zap.Logger
}

func NewGCSStorage(cfg StorageConfig, logger *zap.Logger) *GCSStorage {
	return &GCSStorage{
		logger: logger,
	}
}

func (g *GCSStorage) GetState(workspace string) ([]byte, error)    { return nil, nil }
func (g *GCSStorage) PutState(workspace string, data []byte) error { return nil }
func (g *GCSStorage) ListModules() ([]string, error)               { return []string{}, nil }
func (g *GCSStorage) UploadModule(name string, data []byte) error  { return nil }
