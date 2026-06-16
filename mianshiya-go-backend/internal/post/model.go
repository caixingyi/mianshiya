package post

import "time"

type Post struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Tags      string    `gorm:"type:json" json:"tags"`
	ThumbNum  int       `gorm:"not null;default:0" json:"thumbNum"`
	FavourNum int       `gorm:"not null;default:0" json:"favourNum"`
	UserID    int64     `gorm:"not null" json:"userId"`
	IsDelete  int       `gorm:"not null;default:0" json:"isDelete"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
