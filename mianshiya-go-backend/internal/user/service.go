package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

const salt = "mianshiya"

func encryptPassword(password string) string {
	sum := md5.Sum([]byte(salt + password))
	return hex.EncodeToString(sum[:])
}
func (s *Service) Register(req *RegisterRequest) (int64, error) {
	// 1. 校验参数
	if req.UserAccount == "" || req.UserPassword == "" || req.CheckPassword == "" {
		return 0, errors.New("账号或密码不能为空")
	}
	if req.UserPassword != req.CheckPassword {
		return 0, errors.New("两次输入的密码不一致")
	}

	if len(req.UserAccount) < 4 {
		return 0, errors.New("账号长度不能少于4位")
	}
	if len(req.UserPassword) < 8 {
		return 0, errors.New("密码长度不能少于8位")
	}

	// 2. 判断账号是否已存在Register
	existingUser, err := s.repo.FindByAccount(req.UserAccount)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if existingUser != nil {
		return 0, errors.New("账号已存在")
	}

	// 3. 加密密码
	encryptedPassword := encryptPassword(req.UserPassword)

	// 4. 构造User
	user := &User{
		UserAccount:  req.UserAccount,
		UserPassword: encryptedPassword,
		UserRole:     "user",
	}

	// 5. 写入数据库
	return s.repo.Create(user)
}
