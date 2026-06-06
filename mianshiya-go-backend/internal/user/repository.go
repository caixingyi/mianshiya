package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByAccount(account string) (*User, error) {
	var user User
	err := r.db.Where("user_account = ? AND is_delete = 0", account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

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
