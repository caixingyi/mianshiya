package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"mianshiya-go-backend/internal/response"

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

// AddUser 管理员添加用户
func (s *Service) AddUser(req *AddUserRequest) (int64, error) {
	// 1. 校验参数
	if req == nil {
		return 0, errors.New("请求参数不能为空")
	}
	if req.UserAccount == "" {
		return 0, errors.New("用户账号不能为空")
	}
	if len(req.UserAccount) < 4 {
		return 0, errors.New("用户账号长度不能少于4位")
	}
	role := req.UserRole
	if role == "" {
		role = UserRoleUser
	}
	if role != UserRoleUser && role != UserRoleAdmin {
		return 0, errors.New("用户角色不合法")
	}
	// 2. 判断账号是否已存在
	existingUser, err := s.repo.FindByAccount(req.UserAccount)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if existingUser != nil {
		return 0, errors.New("账号已存在")
	}
	// 3. 加默认密码
	encryptedPassword := encryptPassword("12345678")
	// 4. 构造 User 对象
	user := &User{
		UserAccount:  req.UserAccount,
		UserName:     req.UserName,
		UserAvatar:   req.UserAvatar,
		UserPassword: encryptedPassword,
		UserRole:     role,
	}
	// 5. 写入数据库
	return s.repo.Create(user)
}

// DeleteUser 管理员删除用户
func (s *Service) DeleteUser(id int64) error {
	if id <= 0 {
		return errors.New("参数错误")
	}
	return s.repo.DeleteByID(id)
}

// UpdateUser 管理员更新用户信息
func (s *Service) UpdateUser(id int64, req *UpdateUserRequest) error {
	if id <= 0 {
		return errors.New("参数错误")
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
	if req.UserRole != "" {
		if req.UserRole != UserRoleUser && req.UserRole != UserRoleAdmin && req.UserRole != UserRoleBan {
			return errors.New("用户角色不合法")
		}
		updates["user_role"] = req.UserRole
	}
	if len(updates) == 0 {
		return errors.New("没有要更新的字段")
	}
	return s.repo.UpdateByID(id, updates)
}

func (s *Service) ListUsers(req *ListUserRequest) (*response.PageResponse[UserResponse], error) {
	// 1. 校验参数
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 50 {
		return nil, errors.New("参数错误")
	}
	if req.UserRole != "" && req.UserRole != UserRoleUser && req.UserRole != UserRoleAdmin && req.UserRole != UserRoleBan {
		return nil, errors.New("用户角色不合法")
	}
	// 2. 调用 Repository 层查询用户列表
	users, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}
	// 3. 构造响应数据
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *toUserResponse(user)
	}
	// 4. 返回分页响应
	return &response.PageResponse[UserResponse]{
		Records:  userResponses,
		Total:    total,
		Current:  req.Current,
		PageSize: req.PageSize,
	}, nil
}

// 转换为用户响应结构
func toUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
	}
}

// GetUserResponseByID 根据用户 ID 获取用户信息
func (s *Service) GetUserResponseByID(id int64) (*UserResponse, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}
