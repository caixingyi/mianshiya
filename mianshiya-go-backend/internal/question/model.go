package question

import "time"

type Question struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"size:256;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Tags      string    `gorm:"type:varchar(1024)" json:"tags"`
	Answer    string    `gorm:"type:text" json:"answer"`
	UserID    int64     `gorm:"not null" json:"userId"`
	EditTime  time.Time `gorm:"autoUpdateTime" json:"editTime"`
	IsDelete  int       `gorm:"not null;default:0" json:"isDelete"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
