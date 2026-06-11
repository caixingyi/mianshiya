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

// Create 创建单条题库与题目的关联关系
func (r *Repository) Create(relation *QuestionBankQuestion) (int64, error) {
	err := r.db.Create(relation).Error
	if err != nil {
		return 0, err
	}
	return relation.ID, nil
}

// FindByID 根据 ID 查询题库与题目的关联关系
func (r *Repository) FindByID(id int64) (*QuestionBankQuestion, error) {
	var relation QuestionBankQuestion
	err := r.db.Where("id = ?", id).First(&relation).Error
	if err != nil {
		return nil, err
	}
	return &relation, nil
}

// FindByBankIDAndQuestionID 根据题库 ID 和题目 ID 查询关联关系
func (r *Repository) FindByBankIDAndQuestionID(questionBankID, questionID int64) (*QuestionBankQuestion, error) {
	var relation QuestionBankQuestion
	err := r.db.Where("question_bank_id = ? AND question_id = ?", questionBankID, questionID).First(&relation).Error
	if err != nil {
		return nil, err
	}
	return &relation, nil
}

// List 分页查询题库与题目的关联关系列表
func (r *Repository) List(req *ListQuestionBankQuestionRequest) ([]*QuestionBankQuestion, int64, error) {
	result := make([]*QuestionBankQuestion, 0)
	query := r.db.Model(&QuestionBankQuestion{})

	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.NotID > 0 {
		query = query.Where("id <> ?", req.NotID)
	}
	if req.QuestionBankID > 0 {
		query = query.Where("question_bank_id = ?", req.QuestionBankID)
	}
	if req.QuestionID > 0 {
		query = query.Where("question_id = ?", req.QuestionID)
	}
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(int((req.Current - 1) * req.PageSize)).
		Limit(int(req.PageSize)).
		Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// UpdateByID 根据 ID 更新题库与题目的关联关系
func (r *Repository) UpdateByID(id int64, updates map[string]any) error {
	result := r.db.Model(&QuestionBankQuestion{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteByID 根据 ID 删除题库与题目的关联关系
func (r *Repository) DeleteByID(id int64) error {
	result := r.db.Where("id = ?", id).Delete(&QuestionBankQuestion{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteByBankIDAndQuestionID 根据题库 ID 和题目 ID 删除关联关系
func (r *Repository) DeleteByBankIDAndQuestionID(questionBankID, questionID int64) error {
	result := r.db.Where("question_bank_id = ? AND question_id = ?", questionBankID, questionID).
		Delete(&QuestionBankQuestion{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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
