package post

import (
	"mianshiya-go-backend/internal/user"
	"time"
)

// PostResponse 定义了帖子响应的结构体
type PostResponse struct {
	//	帖子相关字段
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"userId"`
	TagList   []string  `json:"tagList"`
	ThumbNum  int       `json:"thumbNum"`
	FavourNum int       `json:"favourNum"`
	HasThumb  bool      `json:"hasThumb"`
	HasFavour bool      `json:"hasFavour"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// 用户相关信息字段
	User *user.UserResponse `json:"user"`
}
