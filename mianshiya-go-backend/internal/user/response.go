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
