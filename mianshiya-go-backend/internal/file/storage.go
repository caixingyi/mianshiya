package file

import (
	"context"
	"io"
)

// Storage 文件存储接口
type Storage interface {
	Save(ctx context.Context, objectKey string, reader io.Reader) (string, error)
}
