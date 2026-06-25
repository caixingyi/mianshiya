package mockinterview

import "gorm.io/gorm"

// Repository 定义了模拟面试相关的数据访问层
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建一个新的 Repository 实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建模拟面试记录
func (r *Repository) Create(mockInterview *MockInterview) (int64, error) {
	err := r.db.Create(mockInterview).Error
	if err != nil {
		return 0, err
	}
	return mockInterview.ID, nil
}

// FindByID 根据 ID 查询模拟面试记录
func (r *Repository) FindByID(id int64) (*MockInterview, error) {
	var mockInterview MockInterview
	err := r.db.Where("id = ? AND is_delete = 0", id).First(&mockInterview).Error
	if err != nil {
		return nil, err
	}
	return &mockInterview, nil
}

// List 分页查询模拟面试记录
func (r *Repository) List(req *ListMockInterviewRequest) ([]*MockInterview, int64, error) {
	result := make([]*MockInterview, 0)
	query := r.db.Model(&MockInterview{}).Where("is_delete = 0")

	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.WorkExperience != "" {
		query = query.Where("work_experience LIKE ?", "%"+req.WorkExperience+"%")
	}
	if req.JobPosition != "" {
		query = query.Where("job_position LIKE ?", "%"+req.JobPosition+"%")
	}
	if req.Difficulty != "" {
		query = query.Where("difficulty LIKE ?", "%"+req.Difficulty+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(int((req.Current - 1) * req.PageSize)).Limit(int(req.PageSize)).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

// UpdateByID 根据 ID 更新模拟面试记录
func (r *Repository) UpdateByID(id int64, updates map[string]any) error {
	result := r.db.Model(&MockInterview{}).
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

// DeleteByID 根据 ID 删除模拟面试记录（软删除）
func (r *Repository) DeleteByID(id int64) error {
	result := r.db.Model(&MockInterview{}).
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
