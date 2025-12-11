package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type LocalStorage struct {
	basePath      string
	basePathTrash string
}

func NewLocalStorage(basePath, basePathTrash string) *LocalStorage {
	return &LocalStorage{basePath: basePath, basePathTrash: basePathTrash}
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
	return path, err
}

func (localStorage *LocalStorage) MoveFile(originalPath, targetFolder string) (string, error) {
	srcPath := filepath.Join(localStorage.basePath, originalPath)

	// Ensure target folder exists
	targetDir := filepath.Join(localStorage.basePathTrash, targetFolder)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", err
	}

	// Create new filename with timestamp to avoid conflicts
	filename := filepath.Base(originalPath)
	newFilename := time.Now().Format("20060102_150405") + "_" + filename
	destPath := filepath.Join(targetDir, newFilename)

	// Move
	if err := os.Rename(srcPath, destPath); err != nil {
		return "", err
	}

	// Return path relative to base
	return filepath.Join(targetFolder, newFilename), nil
}

func (localStorage *LocalStorage) Delete(path string) error {
	return os.Remove(filepath.Join(localStorage.basePath, path))
}
