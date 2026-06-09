package questionbankquestion

import "gorm.io/gorm"

// Repository 定义了题库与题目的关联关系的存储库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建一个新的 Repository 实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// BatchCreate 批量创建题库与题目的关联关系
func (r *Repository) BatchCreate(relations []QuestionBankQuestion) error {
	if len(relations) == 0 {
		return nil
	}
	return r.db.Create(&relations).Error
}

// FindExistingQuestionIDs 查找题库中已经存在的题目ID列表
func (r *Repository) FindExistingQuestionIDs(questionBankID int64, questionIDs []int64) ([]int64, error) {
	existingIDs := make([]int64, 0)

	if len(questionIDs) == 0 {
		return existingIDs, nil
	}
	if err := r.db.Model(&QuestionBankQuestion{}).
		Where("question_bank_id = ? AND question_id IN ?", questionBankID, questionIDs).
		Pluck("question_id", &existingIDs).
		Error; err != nil {
		return nil, err
	}
	return existingIDs, nil
}

// BatchRemove 批量删除题库与题目的关联关系
func (r *Repository) BatchRemove(questionBankID int64, questionIDs []int64) (int64, error) {
	if len(questionIDs) == 0 {
		return 0, nil
	}
	result := r.db.Where("question_bank_id = ? AND question_id IN ?", questionBankID, questionIDs).
		Delete(&QuestionBankQuestion{})
	return result.RowsAffected, result.Error
}

// FindQuestionIDsByBankID 根据题库ID查找关联的题目ID列表
func (r *Repository) FindQuestionIDsByBankID(questionBankID int64) ([]int64, error) {
	questionIDs := make([]int64, 0)
	if err := r.db.Model(&QuestionBankQuestion{}).
		Where("question_bank_id = ?", questionBankID).
		Pluck("question_id", &questionIDs).
		Error; err != nil {
		return nil, err
	}
	return questionIDs, nil
}
