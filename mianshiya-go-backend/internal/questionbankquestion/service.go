package questionbankquestion

import (
	"errors"
	"mianshiya-go-backend/internal/question"
	"mianshiya-go-backend/internal/questionbank"
)

// Service 定义了题库与题目的关联关系的服务层
type Service struct {
	repo             *Repository
	questionRepo     *question.Repository
	questionBankRepo *questionbank.Repository
}

// NewService 创建一个新的 Service 实例
func NewService(repo *Repository, questionRepo *question.Repository, questionBankRepo *questionbank.Repository) *Service {
	return &Service{repo: repo, questionRepo: questionRepo, questionBankRepo: questionBankRepo}
}

// 辅助函数dedupeIDs用于去重和过滤无效ID，私有函数，外部不可见
func dedupeIDs(ids []int64) []int64 {
	seen := make(map[int64]bool)
	result := make([]int64, 0)
	for _, id := range ids {
		if id <= 0 || seen[id] {
			continue
		}
		seen[id] = true
		result = append(result, id)
	}
	return result
}

// BatchAddQuestionsToBank 批量添加题目到题库
func (s *Service) BatchAddQuestionsToBank(req *BatchAddQuestionsToBankRequest, userID int64) error {
	// 参数校验
	if req == nil {
		return errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return errors.New("用户ID无效")
	}
	if req.QuestionBankID <= 0 {
		return errors.New("题库ID无效")
	}

	questionIDs := dedupeIDs(req.QuestionIDs)
	if len(questionIDs) == 0 {
		return errors.New("题目ID列表不能为空")
	}

	// 检查题库是否存在
	if _, err := s.questionBankRepo.FindByID(req.QuestionBankID); err != nil {
		return errors.New("题库不存在")
	}

	// 检查题目是否存在
	for _, questionID := range questionIDs {
		if _, err := s.questionRepo.FindByID(questionID); err != nil {
			return errors.New("题目不存在")
		}
	}

	// 查已经存在的题目 ID
	existingIDs, err := s.repo.FindExistingQuestionIDs(req.QuestionBankID, questionIDs)
	if err != nil {
		return err
	}

	// 把已有 ID 转成 map，方便判断
	existingMap := make(map[int64]bool)
	for _, id := range existingIDs {
		existingMap[id] = true
	}

	relations := make([]QuestionBankQuestion, 0)
	for _, questionID := range questionIDs {
		if existingMap[questionID] {
			continue
		}

		relations = append(relations, QuestionBankQuestion{
			QuestionBankID: req.QuestionBankID,
			QuestionID:     questionID,
			UserID:         userID,
		})
	}

	if len(relations) == 0 {
		return errors.New("题目已全部存在于题库中")
	}

	return s.repo.BatchCreate(relations)
}

func (s *Service) BatchRemoveQuestionsFromBank(req *BatchRemoveQuestionsFromBankRequest) error {
	// 参数校验
	if req == nil {
		return errors.New("请求参数不能为空")
	}
	if req.QuestionBankID <= 0 {
		return errors.New("题库ID无效")
	}
	// 题目ID列表去重和过滤无效ID
	questionIDs := dedupeIDs(req.QuestionIDs)
	if len(questionIDs) == 0 {
		return errors.New("题目ID列表不能为空")
	}
	// 检查题库是否存在
	if _, err := s.questionBankRepo.FindByID(req.QuestionBankID); err != nil {
		return errors.New("题库不存在")
	}
	rowsAffected, err := s.repo.BatchRemove(req.QuestionBankID, questionIDs)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("没有可移除的题目")
	}

	return nil
}
