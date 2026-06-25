package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	baseDir string
	urlBase string
}

// NewLocalStorage 创建本地文件存储
func NewLocalStorage(baseDir string, urlBase string) *LocalStorage {
	return &LocalStorage{
		baseDir: baseDir,
		urlBase: urlBase,
	}
}

// Save 保存文件到本地磁盘，并返回可访问 URL
func (s *LocalStorage) Save(ctx context.Context, objectKey string, reader io.Reader) (string, error) {
	fullPath := filepath.Join(s.baseDir, objectKey)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return "", err
	}

	return s.urlBase + "/" + filepath.ToSlash(objectKey), nil
}
