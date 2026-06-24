package postthumb

// AddPostThumbRequest 点赞/取消点赞帖子的请求参数
type AddPostThumbRequest struct {
	PostID int64 `json:"postId" binding:"required"`
}
