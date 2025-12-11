package storage

import "mime/multipart"

type Contract interface {
	Put(file *multipart.FileHeader, destPath string) (string, error)
	MoveFile(originalPath, targetFolder string) (string, error)
	Delete(path string) error
}
