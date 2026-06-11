package questionbankquestion

import (
	"errors"
	"mianshiya-go-backend/internal/question"
	"mianshiya-go-backend/internal/questionbank"
	"mianshiya-go-backend/internal/response"

	"gorm.io/gorm"
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

func (s *Service) toQuestionBankQuestionResponse(relation *QuestionBankQuestion) *QuestionBankQuestionResponse {
	return &QuestionBankQuestionResponse{
		ID:             relation.ID,
		QuestionBankID: relation.QuestionBankID,
		QuestionID:     relation.QuestionID,
		UserID:         relation.UserID,
		CreatedAt:      relation.CreatedAt,
		UpdatedAt:      relation.UpdatedAt,
	}
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

func (s *Service) validateQuestionBankAndQuestion(questionBankID, questionID int64) error {
	if questionBankID <= 0 {
		return errors.New("题库ID无效")
	}
	if questionID <= 0 {
		return errors.New("题目ID无效")
	}
	if _, err := s.questionBankRepo.FindByID(questionBankID); err != nil {
		return errors.New("题库不存在")
	}
	if _, err := s.questionRepo.FindByID(questionID); err != nil {
		return errors.New("题目不存在")
	}
	return nil
}

// AddQuestionBankQuestion 添加单条题库题目关联
func (s *Service) AddQuestionBankQuestion(req *AddQuestionBankQuestionRequest, userID int64) (int64, error) {
	if req == nil {
		return 0, errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return 0, errors.New("用户ID无效")
	}
	if err := s.validateQuestionBankAndQuestion(req.QuestionBankID, req.QuestionID); err != nil {
		return 0, err
	}

	_, err := s.repo.FindByBankIDAndQuestionID(req.QuestionBankID, req.QuestionID)
	if err == nil {
		return 0, errors.New("题目已存在于题库中")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	relation := &QuestionBankQuestion{
		QuestionBankID: req.QuestionBankID,
		QuestionID:     req.QuestionID,
		UserID:         userID,
	}
	return s.repo.Create(relation)
}

// DeleteQuestionBankQuestion 删除单条题库题目关联
func (s *Service) DeleteQuestionBankQuestion(id int64) error {
	if id <= 0 {
		return errors.New("参数错误")
	}
	return s.repo.DeleteByID(id)
}

// UpdateQuestionBankQuestion 更新题库题目关联
func (s *Service) UpdateQuestionBankQuestion(req *UpdateQuestionBankQuestionRequest) error {
	if req == nil || req.ID <= 0 {
		return errors.New("参数错误")
	}

	oldRelation, err := s.repo.FindByID(req.ID)
	if err != nil {
		return err
	}

	questionBankID := oldRelation.QuestionBankID
	if req.QuestionBankID > 0 {
		questionBankID = req.QuestionBankID
	}

	questionID := oldRelation.QuestionID
	if req.QuestionID > 0 {
		questionID = req.QuestionID
	}

	if err := s.validateQuestionBankAndQuestion(questionBankID, questionID); err != nil {
		return err
	}

	existingRelation, err := s.repo.FindByBankIDAndQuestionID(questionBankID, questionID)
	if err == nil && existingRelation.ID != req.ID {
		return errors.New("题目已存在于题库中")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	updates := make(map[string]any)
	if req.QuestionBankID > 0 {
		updates["question_bank_id"] = req.QuestionBankID
	}
	if req.QuestionID > 0 {
		updates["question_id"] = req.QuestionID
	}
	if len(updates) == 0 {
		return errors.New("没有要更新的字段")
	}

	return s.repo.UpdateByID(req.ID, updates)
}

// GetQuestionBankQuestionResponseByID 根据 ID 获取题库题目关联详情
func (s *Service) GetQuestionBankQuestionResponseByID(id int64) (*QuestionBankQuestionResponse, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	relation, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toQuestionBankQuestionResponse(relation), nil
}

// ListQuestionBankQuestions 分页获取题库题目关联 VO 列表
func (s *Service) ListQuestionBankQuestions(req *ListQuestionBankQuestionRequest) (*response.PageResponse[QuestionBankQuestionResponse], error) {
	if err := normalizeListRequest(req); err != nil {
		return nil, err
	}

	records, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}

	responses := make([]QuestionBankQuestionResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, *s.toQuestionBankQuestionResponse(record))
	}

	return &response.PageResponse[QuestionBankQuestionResponse]{
		Records:  responses,
		Total:    total,
		Current:  req.Current,
		PageSize: req.PageSize,
	}, nil
}

// ListQuestionBankQuestionPage 分页获取题库题目关联原始列表
func (s *Service) ListQuestionBankQuestionPage(req *ListQuestionBankQuestionRequest) (*response.PageResponse[QuestionBankQuestion], error) {
	if err := normalizeListRequest(req); err != nil {
		return nil, err
	}

	records, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}

	relations := make([]QuestionBankQuestion, 0, len(records))
	for _, record := range records {
		relations = append(relations, *record)
	}

	return &response.PageResponse[QuestionBankQuestion]{
		Records:  relations,
		Total:    total,
		Current:  req.Current,
		PageSize: req.PageSize,
	}, nil
}

func normalizeListRequest(req *ListQuestionBankQuestionRequest) error {
	if req == nil {
		return errors.New("请求参数不能为空")
	}
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 20 {
		return errors.New("参数错误")
	}
	return nil
}

// RemoveQuestionBankQuestion 按题库 ID 和题目 ID 移除关联
func (s *Service) RemoveQuestionBankQuestion(req *RemoveQuestionBankQuestionRequest) error {
	if req == nil {
		return errors.New("请求参数不能为空")
	}
	if req.QuestionBankID <= 0 || req.QuestionID <= 0 {
		return errors.New("参数错误")
	}
	return s.repo.DeleteByBankIDAndQuestionID(req.QuestionBankID, req.QuestionID)
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
