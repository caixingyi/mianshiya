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
		UserRole:     UserRoleUser,
	}

	// 5. 写入数据库
	return s.repo.Create(user)
}

func toLoginUserResponse(user *User) *LoginUserResponse {
	return &LoginUserResponse{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserRole:    user.UserRole,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
	}
}

func (s *Service) Login(req *LoginRequest) (*LoginUserResponse, error) {
	// 1. 校验参数
	if req.UserAccount == "" || req.UserPassword == "" {
		return nil, errors.New("账号或密码不能为空")
	}
	if len(req.UserAccount) < 4 {
		return nil, errors.New("账号长度不能少于4位")
	}
	if len(req.UserPassword) < 8 {
		return nil, errors.New("密码长度不能少于8位")
	}

	// 2. 查找用户
	user, err := s.repo.FindByAccount(req.UserAccount)

	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 3. 验证密码
	if user.UserPassword != encryptPassword(req.UserPassword) {
		return nil, errors.New("密码错误")
	}

	// 4. 返回登录成功响应
	return toLoginUserResponse(user), nil
}

func (s *Service) GetUserByID(id int64) (*LoginUserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return toLoginUserResponse(user), nil
}

func (s *Service) UpdateMy(userID int64, req *UpdateMyRequest) error {
	if userID <= 0 {
		return errors.New("无效的用户ID")
	}

	if req == nil {
		return errors.New("请求参数不能为空")
	}

	updates := make(map[string]any)
	if req.UserName != "" {
		updates["user_name"] = req.UserName
	}
	if req.UserAvatar != "" {
		updates["user_avatar"] = req.UserAvatar
	}
	if req.UserProfile != "" {
		updates["user_profile"] = req.UserProfile
	}

	if len(updates) == 0 {
		return errors.New("没有要更新的字段")
	}

	return s.repo.UpdateByID(userID, updates)
}

// IsAdmin 判断用户是否为管理员
func (s *Service) IsAdmin(userID int64) (bool, error) {
	if userID <= 0 {
		return false, errors.New("无效的用户ID")
	}
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return false, err
	}
	return user.UserRole == UserRoleAdmin, nil
}
