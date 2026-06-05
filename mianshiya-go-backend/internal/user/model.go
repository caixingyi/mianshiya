package user

import "time"

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserAccount  string    `gorm:"size:256;not null;uniqueIndex" json:"userAccount"`
	UserPassword string    `gorm:"size:512;not null" json:"-"`
	UserRole     string    `gorm:"size:256;not null;default:user" json:"userRole"`
	IsDelete     int       `gorm:"not null;default:0" json:"isDelete"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
