package questionbankquestion

import "time"

type QuestionBankQuestion struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionBankID int64     `gorm:"not null;uniqueIndex:idx_bank_question" json:"questionBankId"`
	QuestionID     int64     `gorm:"not null;uniqueIndex:idx_bank_question" json:"questionId"`
	UserID         int64     `gorm:"not null;index" json:"userId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
