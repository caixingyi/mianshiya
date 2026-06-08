package questionbank

import (
	"errors"
	"mianshiya-go-backend/internal/response"
)

// Service 定义了题库服务的结构体
type Service struct {
	repo *Repository
}

// NewService 创建一个新的 Service 实例
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) toQuestionBankResponse(qb *QuestionBank) *QuestionBankResponse {
	return &QuestionBankResponse{
		ID:          qb.ID,
		Title:       qb.Title,
		Description: qb.Description,
		Picture:     qb.Picture,
		UserID:      qb.UserID,
		EditTime:    qb.EditTime,
		CreatedAt:   qb.CreatedAt,
		UpdatedAt:   qb.UpdatedAt,
	}
}

// AddQuestionBank 添加新题库
func (s *Service) AddQuestionBank(req *AddQuestionBankRequest, userID int64) (int64, error) {
	// 1. 校验参数
	if req == nil {
		return 0, errors.New("请求参数不能为空")
	}
	if req.Title == "" {
		return 0, errors.New("标题不能为空")
	}
	if userID <= 0 {
		return 0, errors.New("用户ID无效")
	}
	// 2. 创建题库对象
	questionBank := &QuestionBank{
		Title:       req.Title,
		Description: req.Description,
		Picture:     req.Picture,
		UserID:      userID,
	}
	// 3. 调用 Repository 添加题库
	return s.repo.Create(questionBank)
}

// GetQuestionBankResponseByID 根据 ID 获取题库详情
func (s *Service) GetQuestionBankResponseByID(id int64) (*QuestionBankResponse, error) {
	// 1. 校验参数
	if id <= 0 {
		return nil, errors.New("题库ID无效")
	}
	// 2. 调用 Repository 获取题库详情
	questionBank, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	// 3. 转换为响应对象并返回
	return s.toQuestionBankResponse(questionBank), nil
}

// ListQuestionBanks 获取题库列表
func (s *Service) ListQuestionBanks(req *ListQuestionBankRequest) (*response.PageResponse[QuestionBankResponse], error) {
	// 1. 校验参数
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		return nil, errors.New("每页记录数不能超过100")
	}
	questionBanks, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}
	// 3. 转换为响应对象并返回
	records := make([]QuestionBankResponse, len(questionBanks))
	for i, qb := range questionBanks {
		records[i] = *s.toQuestionBankResponse(qb)
	}
	return &response.PageResponse[QuestionBankResponse]{
		Records:  records,
		Total:    total,
		Current:  req.Current,
		PageSize: req.PageSize,
	}, nil
}
