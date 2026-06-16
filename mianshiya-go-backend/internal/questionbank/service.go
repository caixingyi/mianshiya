package questionbank

import (
	"errors"
	"mianshiya-go-backend/internal/question"
	"mianshiya-go-backend/internal/response"
)

// Service 定义了题库服务的结构体
type Service struct {
	repo        *Repository
	questionSvc *question.Service
}

// NewService 创建一个新的 Service 实例
func NewService(repo *Repository, questionSvc *question.Service) *Service {
	return &Service{repo: repo, questionSvc: questionSvc}
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

func (s *Service) GetQuestionBankResponseByID(req *GetQuestionBankRequest) (*QuestionBankResponse, error) {
	// 1. 校验参数
	if req.ID <= 0 {
		return nil, errors.New("题库ID无效")
	}
	// 2. 调用 Repository 获取题库详情
	questionBank, err := s.repo.FindByID(req.ID)
	if err != nil {
		return nil, err
	}
	// 3. 转换为响应对象
	resp := s.toQuestionBankResponse(questionBank)
	// 4. 如果需要查询题目列表，则通过 question.Service 填充题目分页数据
	if req.NeedQueryQuestionList {
		current := req.Current
		if current <= 0 {
			current = 1
		}
		pageSize := req.PageSize
		if pageSize <= 0 {
			pageSize = 10
		}
		questionPage, err := s.questionSvc.ListQuestions(&question.ListQuestionRequest{
			QuestionBankID: req.ID,
			Current:        current,
			PageSize:       pageSize,
		})
		if err != nil {
			return nil, err
		}
		resp.QuestionPage = questionPage
	}
	return resp, nil
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

// DeleteQuestionBank 删除题库
func (s *Service) DeleteQuestionBank(id int64) error {
	// 1. 校验参数
	if id <= 0 {
		return errors.New("参数错误")
	}
	// 2. 调用 Repository 删除题库
	return s.repo.Delete(id)
}

// UpdateQuestionBank 更新题库
func (s *Service) UpdateQuestionBank(req *UpdateQuestionBankRequest) error {
	if req == nil || req.ID <= 0 {
		return errors.New("参数错误")
	}

	updates := make(map[string]any)

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Picture != "" {
		updates["picture"] = req.Picture
	}

	if len(updates) == 0 {
		return errors.New("没有要更新的字段")
	}

	return s.repo.UpdateByID(req.ID, updates)
}

// ListQuestionBankPage 获取题库列表（分页）
func (s *Service) ListQuestionBankPage(req *ListQuestionBankRequest) (*response.PageResponse[QuestionBank], error) {
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

	records := make([]QuestionBank, 0, len(questionBanks))
	for _, qb := range questionBanks {
		records = append(records, *qb)
	}

	return &response.PageResponse[QuestionBank]{
		Records:  records,
		Total:    total,
		Current:  req.Current,
		PageSize: req.PageSize,
	}, nil
}
