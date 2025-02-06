package storage

type GCSStorage struct{}

func NewGCSStorage(cfg StorageConfig) *GCSStorage {
	return &GCSStorage{}
}

func (g *GCSStorage) GetState(workspace string) ([]byte, error)    { return nil, nil }
func (g *GCSStorage) PutState(workspace string, data []byte) error { return nil }
func (g *GCSStorage) ListModules() ([]string, error)               { return []string{}, nil }
func (g *GCSStorage) UploadModule(name string, data []byte) error  { return nil }
