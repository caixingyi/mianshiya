package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByAccount 根据账号查找用户
func (r *Repository) FindByAccount(account string) (*User, error) {
	var user User
	err := r.db.Where("user_account = ? AND is_delete = 0", account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建新用户
func (r *Repository) Create(user *User) (int64, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// 按照 ID 查找用户
func (r *Repository) FindByID(id int64) (*User, error) {
	var user User
	err := r.db.Where("id = ? AND is_delete = 0", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateByID 根据 ID 更新用户信息
func (r *Repository) UpdateByID(id int64, updates map[string]any) error {
	result := r.db.Model(&User{}).Where("id = ? AND is_delete = 0", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	// 如果没有任何记录被更新，说明用户不存在或已经被删除
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteByID 根据 ID 删除用户（软删除）
func (r *Repository) DeleteByID(id int64) error {
	result := r.db.Model(&User{}).Where("id = ? AND is_delete = 0", id).Update("is_delete", 1)
	if result.Error != nil {
		return result.Error
	}
	// 如果没有任何记录被更新，说明用户不存在或已经被删除
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
