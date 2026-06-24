package postthumb

import "gorm.io/gorm"

// Repository 定义了帖子点赞相关的数据访问层
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建一个新的 Repository 实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByPostIDAndUserID 根据帖子 ID 和用户 ID 查询点赞记录
func (r *Repository) FindByPostIDAndUserID(postID, userID int64) (*PostThumb, error) {
	var postThumb PostThumb
	err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&postThumb).Error
	if err != nil {
		return nil, err
	}
	return &postThumb, nil
}

// Create 创建点赞记录
func (r *Repository) Create(postThumb *PostThumb) error {
	return r.db.Create(postThumb).Error
}

// DeleteByPostIDAndUserID 根据帖子 ID 和用户 ID 删除点赞记录
func (r *Repository) DeleteByPostIDAndUserID(postID, userID int64) error {
	result := r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&PostThumb{})
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
