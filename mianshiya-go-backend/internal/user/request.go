package user

// 用户注册请求参数
type RegisterRequest struct {
	UserAccount   string `json:"userAccount"`
	UserPassword  string `json:"userPassword"`
	CheckPassword string `json:"checkPassword"`
}

// 登录请求参数
type LoginRequest struct {
	UserAccount  string `json:"userAccount"`
	UserPassword string `json:"userPassword"`
}

// 更新用户信息请求参数
type UpdateMyRequest struct {
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	UserProfile string `json:"userProfile"`
}

// 编辑当前登录用户信息请求参数
type EditUserRequest struct {
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	UserProfile string `json:"userProfile"`
}

// 管理员添加用户请求参数
type AddUserRequest struct {
	UserAccount string `json:"userAccount"`
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	UserRole    string `json:"userRole"`
}

// 删除用户请求参数
type DeleteUserRequest struct {
	ID int64 `json:"id"`
}

// 管理员更新用户信息请求参数
type UpdateUserRequest struct {
	ID          int64  `json:"id"`
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	UserProfile string `json:"userProfile"`
	UserRole    string `json:"userRole"`
}

// 列表用户请求参数
type ListUserRequest struct {
	Current  int64  `json:"current"`
	PageSize int64  `json:"pageSize"`
	ID       int64  `json:"id"`
	UserName string `json:"userName"`
	UserRole string `json:"userRole"`
}

// 根据 ID 获取用户信息请求参数
type GetUserRequest struct {
	ID int64 `form:"id"`
}

// GetUserSignInRequest 获取用户签到记录请求参数
type GetUserSignInRequest struct {
	Year int `form:"year"`
}
