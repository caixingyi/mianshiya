package post

// AddPostRequest 添加帖子的请求参数
type AddPostRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

// UpdatePostRequest 更新帖子的请求参数
type UpdatePostRequest struct {
	ID      int64    `json:"id" binding:"required"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// EditPostRequest 编辑帖子的请求参数
type EditPostRequest struct {
	ID      int64    `json:"id" binding:"required"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// DeletePostRequest 删除帖子的请求参数
type DeletePostRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// ListPostsRequest 获取帖子列表的请求参数
type ListPostsRequest struct {
	Current      int64    `json:"current" binding:"required"`
	PageSize     int64    `json:"pageSize" binding:"required"`
	ID           int64    `json:"id"`
	NotID        int64    `json:"notId"`
	SearchText   string   `json:"searchText"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Tags         []string `json:"tags"`
	UserID       int64    `json:"userId"`
	FavourUserID int64    `json:"favourUserId"`
}

// GetPostRequest 获取帖子详情的请求参数
type GetPostRequest struct {
	ID int64 `form:"id" binding:"required"`
}
