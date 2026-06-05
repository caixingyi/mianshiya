package user

type LoginUserResponse struct {
	ID          int64  `json:"id"`
	UserAccount string `json:"userAccount"`
	UserRole    string `json:"userRole"`
}
