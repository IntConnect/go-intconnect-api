package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(base string) *LocalStorage {
	return &LocalStorage{basePath: base}
}

func (localStorage *LocalStorage) Put(fileHeader *multipart.FileHeader, path string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fullPath := filepath.Join(localStorage.basePath, path)
	os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return fullPath, err
}

func (localStorage *LocalStorage) Delete(path string) error {
	return os.Remove(filepath.Join(localStorage.basePath, path))
}
