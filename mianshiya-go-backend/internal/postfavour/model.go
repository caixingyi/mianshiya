package postfavour

import "time"

// PostFavour 表示用户对帖子的收藏记录
type PostFavour struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int64     `gorm:"not null;uniqueIndex:idx_post_favour_user" json:"postId"`
	UserID    int64     `gorm:"not null;uniqueIndex:idx_post_favour_user" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
