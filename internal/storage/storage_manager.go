package storage

import (
	"errors"
	"go-intconnect-api/configs"
)

type Manager struct {
	diskContract Contract
}

func NewStorageManager(storageConfig configs.StorageConfig) (*Manager, error) {
	var diskContract Contract

	switch storageConfig.Driver {
	case "local":
		diskContract = NewLocalStorage(storageConfig.BasePath)
	default:
		return nil, errors.New("unsupported storage driver: " + storageConfig.Driver)
	}

	return &Manager{diskContract: diskContract}, nil
}

func (storageManager *Manager) Disk() Contract {
	return storageManager.diskContract
}
