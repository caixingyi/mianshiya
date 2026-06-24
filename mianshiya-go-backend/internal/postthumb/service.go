package postthumb

import (
	"errors"
	"mianshiya-go-backend/internal/post"

	"gorm.io/gorm"
)

// Service 帖子点赞服务层
type Service struct {
	repo     *Repository
	postRepo *post.Repository
}

// NewService 创建帖子点赞服务实例
func NewService(repo *Repository, postRepo *post.Repository) *Service {
	return &Service{repo: repo, postRepo: postRepo}
}

// DoPostThumb 点赞 / 取消点赞
func (s *Service) DoPostThumb(postID int64, userID int64) (int, error) {
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

	// 2. 点赞记录和点赞数更新必须在同一个事务里
	err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		txThumbRepo := s.repo.WithTx(tx)
		txPostRepo := s.postRepo.WithTx(tx)

		oldThumb, err := txThumbRepo.FindByPostIDAndUserID(postID, userID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 3. 已点赞：取消点赞
		if oldThumb != nil {
			if err := txThumbRepo.DeleteByPostIDAndUserID(postID, userID); err != nil {
				return err
			}
			if err := txPostRepo.IncrementThumbNum(postID, -1); err != nil {
				return err
			}
			resultNum = -1
			return nil
		}

		// 4. 未点赞：新增点赞
		postThumb := &PostThumb{
			PostID: postID,
			UserID: userID,
		}
		if err := txThumbRepo.Create(postThumb); err != nil {
			return err
		}
		if err := txPostRepo.IncrementThumbNum(postID, 1); err != nil {
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
