package question

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建新题目
func (r *Repository) Create(question *Question) (int64, error) {
	err := r.db.Create(question).Error
	if err != nil {
		return 0, err
	}
	return question.ID, nil
}

// FindByID 根据 ID 查找题目
func (r *Repository) FindByID(id int64) (*Question, error) {
	var question Question
	err := r.db.Where("id = ? AND is_delete = 0", id).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

// List 分页查询题目列表
func (r *Repository) List(req *ListQuestionRequest) ([]*Question, int64, error) {
	result := make([]*Question, 0)
	query := r.db.Model(&Question{}).Where("is_delete = 0")
	// 构建查询条件
	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}

	if req.NotID > 0 {
		query = query.Where("id <> ?", req.NotID)
	}

	if req.SearchText != "" {
		like := "%" + req.SearchText + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR answer LIKE ?", like, like, like)
	}

	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}

	if req.Content != "" {
		query = query.Where("content LIKE ?", "%"+req.Content+"%")
	}

	if req.Answer != "" {
		query = query.Where("answer LIKE ?", "%"+req.Answer+"%")
	}

	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.QuestionBankID > 0 {
		questionIDs := make([]int64, 0)

		err := r.db.Table("question_bank_questions").
			Where("question_bank_id = ?", req.QuestionBankID).
			Pluck("question_id", &questionIDs).Error
		if err != nil {
			return nil, 0, err
		}

		if len(questionIDs) == 0 {
			return result, 0, nil
		}

		query = query.Where("id IN ?", questionIDs)
	}
	for _, tag := range req.Tags {
		if tag != "" {
			query = query.Where("tags LIKE ?", "%"+tag+"%")
		}
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
