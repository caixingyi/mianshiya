package user

// LoginUserResponse 包含登录成功后返回的用户信息
type LoginUserResponse struct {
	ID          int64  `json:"id"`
	UserAccount string `json:"userAccount"`
	UserRole    string `json:"userRole"`
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	UserProfile string `json:"userProfile"`
}

// LoginResponse 包含登录成功后返回的 token 和用户信息
type LoginResponse struct {
	Token string             `json:"token"`
	User  *LoginUserResponse `json:"user"`
}

// PageResponse 包含分页查询时返回的响应结构
type PageResponse[T any] struct {
	Records  []T   `json:"records"`
	Total    int64 `json:"total"`
	Current  int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
}

// UserResponse 包含用户信息的响应结构
type UserResponse struct {
	ID          int64  `json:"id"`
	UserAccount string `json:"userAccount"`
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	UserProfile string `json:"userProfile"`
	UserRole    string `json:"userRole"`
}
