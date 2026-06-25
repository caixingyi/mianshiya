package file

import (
	"context"
	"crypto/rand"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	BizUserAvatar = "user_avatar"

	maxAvatarSize = 1 * 1024 * 1024
)

var allowedAvatarExts = map[string]bool{
	"jpeg": true,
	"jpg":  true,
	"svg":  true,
	"png":  true,
	"webp": true,
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

// UploadFile 上传文件
func (s *Service) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, biz string, userID int64) (string, error) {
	if userID <= 0 {
		return "", errors.New("无效的用户ID")
	}
	if fileHeader == nil {
		return "", errors.New("文件不能为空")
	}
	if biz == "" {
		return "", errors.New("业务类型不能为空")
	}

	if err := validFile(fileHeader, biz); err != nil {
		return "", err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileName := filepath.Base(fileHeader.Filename)
	randomName, err := randomString(8)
	if err != nil {
		return "", err
	}

	objectKey := filepath.ToSlash(filepath.Join(
		biz,
		int64ToString(userID),
		randomName+"-"+fileName,
	))

	return s.storage.Save(ctx, objectKey, src)
}

func validFile(fileHeader *multipart.FileHeader, biz string) error {
	switch biz {
	case BizUserAvatar:
		if fileHeader.Size > maxAvatarSize {
			return errors.New("文件大小不能超过1M")
		}

		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileHeader.Filename)), ".")
		if !allowedAvatarExts[ext] {
			return errors.New("文件类型错误")
		}

		return nil
	default:
		return errors.New("不支持的业务类型")
	}
}

func randomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = letters[int(bytes[i])%len(letters)]
	}
	return string(bytes), nil
}

func int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}
