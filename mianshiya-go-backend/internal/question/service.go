package question

import (
	"encoding/json"
	"errors"
	"mianshiya-go-backend/internal/response"
)

// Service 题目服务层
type Service struct {
	repo *Repository
}

// NewService 创建题目服务实例
func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// 转换 Question 到 QuestionResponse
func (s *Service) toQuestionResponse(question *Question) (*QuestionResponse, error) {
	tagList := make([]string, 0)
	if question.Tags != "" {
		if err := json.Unmarshal([]byte(question.Tags), &tagList); err != nil {
			return nil, err
		}
	}
	return &QuestionResponse{
		ID:        question.ID,
		Title:     question.Title,
		Content:   question.Content,
		Answer:    question.Answer,
		UserID:    question.UserID,
		TagList:   tagList,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
	}, nil
}

// AddQuestion 添加题目
func (s *Service) AddQuestion(req *AddQuestionRequest, userID int64) (int64, error) {
	// 参数校验
	if req == nil {
		return 0, errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return 0, errors.New("无效的用户ID")
	}
	if req.Title == "" {
		return 0, errors.New("题目标题不能为空")
	}
	if req.Content == "" {
		return 0, errors.New("题目内容不能为空")
	}
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}
	// 将标签列表转换为 JSON 字符串存储
	tagsBytes, err := json.Marshal(tags)
	if err != nil {
		return 0, err
	}
	// 创建题目对象并保存到数据库
	question := &Question{
		Title:   req.Title,
		Content: req.Content,
		Tags:    string(tagsBytes),
		Answer:  req.Answer,
		UserID:  userID,
	}
	return s.repo.Create(question)
}

// GetQuestionResponseByID 根据 ID 获取题目详情
func (s *Service) GetQuestionResponseByID(id int64) (*QuestionResponse, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	question, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toQuestionResponse(question)
}

// ListQuestionResponse 分页获取题目列表
func (s *Service) ListQuestions(req *ListQuestionRequest) (*response.PageResponse[QuestionResponse], error) {
	// 参数校验
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 200 {
		return nil, errors.New("参数错误")
	}
	// 查询题目列表
	records, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}
	// 转换为响应结构
	responses := make([]QuestionResponse, 0, len(records))
	for _, record := range records {
		response, err := s.toQuestionResponse(record)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}
	// 构建分页响应
	return &response.PageResponse[QuestionResponse]{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		Records:  responses,
	}, nil
}

// ListMyQuestions 获取我的题目列表
func (s *Service) ListMyQuestions(req *ListQuestionRequest, userID int64) (*response.PageResponse[QuestionResponse], error) {
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return nil, errors.New("无效的用户ID")
	}
	req.UserID = userID
	return s.ListQuestions(req)
}

// DeleteQuestion 删除题目
func (s *Service) DeleteQuestion(id int64) error {
	if id <= 0 {
		return errors.New("参数错误")
	}
	return s.repo.DeleteByID(id)
}

// UpdateQuestion 更新题目
func (s *Service) UpdateQuestion(req *UpdateQuestionRequest) error {
	if req == nil || req.ID <= 0 {
		return errors.New("参数错误")
	}

	updates := make(map[string]any)

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Answer != "" {
		updates["answer"] = req.Answer
	}
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return err
		}
		updates["tags"] = string(tagsBytes)
	}
	if len(updates) == 0 {
		return errors.New("没有要更新的字段")
	}

	return s.repo.UpdateByID(req.ID, updates)
}

// EditQuestion 编辑题目（用户接口）
func (s *Service) EditQuestion(req *UpdateQuestionRequest, userID int64) error {
	if req == nil || req.ID <= 0 {
		return errors.New("参数错误")
	}
	if userID <= 0 {
		return errors.New("无效的用户ID")
	}
	question, err := s.repo.FindByID(req.ID)
	if err != nil {
		return err
	}
	if question.UserID != userID {
		return errors.New("无权限编辑该题目")
	}
	return s.UpdateQuestion(req)
}

// ListQuestionPage 获取题目分页列表
func (s *Service) ListQuestionPage(req *ListQuestionRequest) (*response.PageResponse[Question], error) {
	// 参数校验
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 20 {
		return nil, errors.New("参数错误")
	}
	// 查询题目列表
	records, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}
	questionList := make([]Question, 0, len(records))
	for _, record := range records {
		questionList = append(questionList, *record)
	}
	return &response.PageResponse[Question]{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		Records:  questionList,
	}, nil
}

// BatchDeleteQuestions 批量删除题目
func (s *Service) BatchDeleteQuestions(req *BatchDeleteQuestionRequest) error {
	if req == nil || len(req.QuestionIDList) == 0 {
		return errors.New("参数错误")
	}
	for _, id := range req.QuestionIDList {
		if id <= 0 {
			return errors.New("参数错误")
		}
	}
	return s.repo.DeleteBatchByIDs(req.QuestionIDList)
}
