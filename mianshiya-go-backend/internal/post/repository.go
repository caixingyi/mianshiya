package post

import "gorm.io/gorm"

// Repository 定义了与 Post 相关的数据访问方法
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建新的 Repository 实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建新帖子
func (r *Repository) Create(post *Post) (int64, error) {
	err := r.db.Create(post).Error
	if err != nil {
		return 0, err
	}
	return post.ID, nil
}

// FindByID 根据 ID 查找帖子
func (r *Repository) FindByID(id int64) (*Post, error) {
	var post Post
	err := r.db.Where("id = ? AND is_delete = 0", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// List 分页查询帖子列表
func (r *Repository) List(req *ListPostsRequest) ([]*Post, int64, error) {
	result := make([]*Post, 0)
	query := r.db.Model(&Post{}).Where("is_delete = 0")
	// 构建查询条件
	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}

	if req.NotID > 0 {
		query = query.Where("id <> ?", req.NotID)
	}

	if req.SearchText != "" {
		like := "%" + req.SearchText + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", like, like)
	}

	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}

	if req.Content != "" {
		query = query.Where("content LIKE ?", "%"+req.Content+"%")
	}
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.FavourUserID > 0 {
		query = query.Where("id IN (SELECT post_id FROM post_favours WHERE user_id = ?)", req.FavourUserID)
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

	if err := query.Offset(int((req.Current - 1) * req.PageSize)).Limit(int(req.PageSize)).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

// UpdateByID 根据 ID 更新帖子
func (r *Repository) UpdateByID(id int64, updates map[string]any) error {
	result := r.db.Model(&Post{}).
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

// DeleteByID 根据 ID 删除帖子（软删除）
func (r *Repository) DeleteByID(id int64) error {
	result := r.db.Model(&Post{}).
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

// DeleteBatchByIDs 根据 ID 列表批量删除帖子（软删除）
func (r *Repository) DeleteBatchByIDs(ids []int64) error {
	result := r.db.Model(&Post{}).
		Where("id IN ? AND is_delete = 0", ids).
		Update("is_delete", 1)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// IncrementThumbNum 原子更新帖子的点赞数
func (r *Repository) IncrementThumbNum(id int64, delta int) error {
	query := r.db.Model(&Post{}).Where("id = ? AND is_delete = 0", id)
	if delta < 0 {
		query = query.Where("thumb_num > 0")
	}
	result := query.UpdateColumn("thumb_num", gorm.Expr("thumb_num + ?", delta))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// IncrementFavourNum 原子更新帖子的收藏数
func (r *Repository) IncrementFavourNum(id int64, delta int) error {
	query := r.db.Model(&Post{}).Where("id = ? AND is_delete = 0", id)
	if delta < 0 {
		query = query.Where("favour_num > 0")
	}
	result := query.UpdateColumn("favour_num", gorm.Expr("favour_num + ?", delta))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// WithTx 返回使用事务的 Repository
func (r *Repository) WithTx(tx *gorm.DB) *Repository {
	return &Repository{db: tx}
}
