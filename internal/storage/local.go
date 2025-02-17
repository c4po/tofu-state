package storage

import (
	"os"

	"go.uber.org/zap"
)

type LocalStorage struct {
	path   string
	logger *zap.Logger
}

func NewLocalStorage(cfg StorageConfig, logger *zap.Logger) *LocalStorage {
	os.MkdirAll(cfg.LocalPath, 0755)
	return &LocalStorage{path: cfg.LocalPath, logger: logger}
}

func (l *LocalStorage) GetState(workspace string) ([]byte, error)    { return nil, nil }
func (l *LocalStorage) PutState(workspace string, data []byte) error { return nil }
func (l *LocalStorage) ListModules() ([]string, error)               { return []string{}, nil }
func (l *LocalStorage) UploadModule(name string, data []byte) error  { return nil }
