package postfavour

import (
	"errors"
	"mianshiya-go-backend/internal/post"

	"gorm.io/gorm"
)

// Service 帖子收藏服务层
type Service struct {
	repo     *Repository
	postRepo *post.Repository
}

// NewService 创建帖子收藏服务实例
func NewService(repo *Repository, postRepo *post.Repository) *Service {
	return &Service{repo: repo, postRepo: postRepo}
}

// DoPostFavour 收藏 / 取消收藏
func (s *Service) DoPostFavour(postID int64, userID int64) (int, error) {
	if postID <= 0 {
		return 0, errors.New("帖子ID无效")
	}
	if userID <= 0 {
		return 0, errors.New("用户ID无效")
	}

	// 1. 先确认帖子存在
	if _, err := s.postRepo.FindByID(postID); err != nil {
		return 0, err
	}

	resultNum := 0

	// 2. 收藏记录和收藏数更新必须在同一个事务里
	err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		txFavourRepo := s.repo.WithTx(tx)
		txPostRepo := s.postRepo.WithTx(tx)

		oldFavour, err := txFavourRepo.FindByPostIDAndUserID(postID, userID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 3. 已收藏：取消收藏
		if oldFavour != nil {
			if err := txFavourRepo.DeleteByPostIDAndUserID(postID, userID); err != nil {
				return err
			}
			if err := txPostRepo.IncrementFavourNum(postID, -1); err != nil {
				return err
			}
			resultNum = -1
			return nil
		}

		// 4. 未收藏：新增收藏
		postFavour := &PostFavour{
			PostID: postID,
			UserID: userID,
		}
		if err := txFavourRepo.Create(postFavour); err != nil {
			return err
		}
		if err := txPostRepo.IncrementFavourNum(postID, 1); err != nil {
			return err
		}
		resultNum = 1
		return nil
	})

	if err != nil {
		return 0, err
	}
	return resultNum, nil
}
