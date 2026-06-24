package postfavour

import "mianshiya-go-backend/internal/post"

// AddPostFavourRequest 收藏/取消收藏帖子的请求参数
type AddPostFavourRequest struct {
	PostID int64 `json:"postId" binding:"required"`
}

// PostFavourQueryRequest 查询用户收藏帖子的请求参数
type PostFavourQueryRequest struct {
	post.ListPostsRequest
	UserID int64 `form:"userId" json:"userId"`
}
