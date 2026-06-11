package questionbank

import "gorm.io/gorm"

// Repository 定义了题库相关的数据访问层
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建一个新的 Repository 实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建新题目
func (r *Repository) Create(questionBank *QuestionBank) (int64, error) {
	result := r.db.Create(questionBank)
	return questionBank.ID, result.Error
}

// FindByID 根据 ID 查找题目
func (r *Repository) FindByID(id int64) (*QuestionBank, error) {
	var question QuestionBank
	err := r.db.Where("id = ? AND is_delete = 0", id).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, err
}

// List 获取题库列表
func (r *Repository) List(req *ListQuestionBankRequest) ([]*QuestionBank, int64, error) {
	result := make([]*QuestionBank, 0)
	var total int64
	// 构建查询条件
	query := r.db.Model(&QuestionBank{}).Where("is_delete = 0")
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Description != "" {
		query = query.Where("description LIKE ?", "%"+req.Description+"%")
	}
	if req.NotID > 0 {
		query = query.Where("id <> ?", req.NotID)
	}
	if req.SearchText != "" {
		like := "%" + req.SearchText + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", like, like)
	}

	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	// 统计总记录数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 分页查询
	err = query.Offset(int((req.Current - 1) * req.PageSize)).Limit(int(req.PageSize)).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// Delete 删除题库
func (r *Repository) Delete(id int64) error {
	result := r.db.Model(&QuestionBank{}).
		Where("id = ? AND is_delete = 0", id).
		Update("is_delete", 1)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateByID 根据 ID 更新题库
func (r *Repository) UpdateByID(id int64, updates map[string]any) error {
	result := r.db.Model(&QuestionBank{}).
		Where("id = ? AND is_delete = 0", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
