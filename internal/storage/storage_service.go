package storage

import "mime/multipart"

type Contract interface {
	Put(file *multipart.FileHeader, destPath string) (string, error)
	Delete(path string) error
}
