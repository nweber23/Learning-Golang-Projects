package storage

import (
	"io"
	"os"
	"path/filepath"
)

type Storage interface {
	Upload(filename string, data io.Reader) (url string, err error)
	Download(filename string) (io.Reader, error)
	Delete(filename string) error
	Exists(filename string) bool
}

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (ls *LocalStorage) Upload(filename string, data io.Reader) (string, error) {
	if err := os.MkdirAll(ls.basePath, os.ModePerm); err != nil {
		return "", err
	}
	filePath := filepath.Join(ls.basePath, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.Copy(file, data); err != nil {
		return "", err
	}
	return filePath, nil
}

func (ls *LocalStorage) Download(filename string) (io.Reader, error) {
	filePath := filepath.Join(ls.basePath, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (ls *LocalStorage) Delete(filename string) error {
	filePath := filepath.Join(ls.basePath, filename)
	return os.Remove(filePath)
}

func (ls *LocalStorage) Exists(filename string) bool {
	filePath := filepath.Join(ls.basePath, filename)
	_, err := os.Stat(filePath)
	return err == nil
}
