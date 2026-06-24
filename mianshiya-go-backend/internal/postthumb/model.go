package postthumb

import "time"

// PostThumb 表示用户对帖子的点赞记录
type PostThumb struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int64     `gorm:"not null;uniqueIndex:idx_post_thumb_user" json:"postId"`
	UserID    int64     `gorm:"not null;uniqueIndex:idx_post_thumb_user" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
