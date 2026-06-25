package mockinterview

import (
	"errors"
	"mianshiya-go-backend/internal/response"
	"mianshiya-go-backend/internal/user"
)

// Service 模拟面试服务层
type Service struct {
	repo    *Repository
	userSvc *user.Service
}

// NewService 创建模拟面试服务实例
func NewService(repo *Repository, userSvc *user.Service) *Service {
	return &Service{repo: repo, userSvc: userSvc}
}

// AddMockInterview 创建模拟面试
func (s *Service) AddMockInterview(req *AddMockInterviewRequest, userID int64) (int64, error) {
	if req == nil {
		return 0, errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return 0, errors.New("无效的用户ID")
	}
	if req.WorkExperience == "" {
		return 0, errors.New("工作年限不能为空")
	}
	if req.JobPosition == "" {
		return 0, errors.New("岗位不能为空")
	}
	if req.Difficulty == "" {
		return 0, errors.New("难度不能为空")
	}

	mockInterview := &MockInterview{
		WorkExperience: req.WorkExperience,
		JobPosition:    req.JobPosition,
		Difficulty:     req.Difficulty,
		Messages:       "",
		Status:         StatusToStart,
		UserID:         userID,
	}

	return s.repo.Create(mockInterview)
}

// DeleteMockInterview 删除模拟面试
func (s *Service) DeleteMockInterview(req *DeleteMockInterviewRequest, userID int64) error {
	if req == nil {
		return errors.New("请求参数不能为空")
	}
	if req.ID <= 0 {
		return errors.New("参数错误")
	}
	if userID <= 0 {
		return errors.New("无效的用户ID")
	}

	mockInterview, err := s.repo.FindByID(req.ID)
	if err != nil {
		return err
	}

	isAdmin, err := s.userSvc.IsAdmin(userID)
	if err != nil {
		return errors.New("无法验证用户权限")
	}
	if mockInterview.UserID != userID && !isAdmin {
		return errors.New("无权限删除该模拟面试")
	}

	return s.repo.DeleteByID(req.ID)
}

// GetMockInterviewByID 根据 ID 获取模拟面试详情
func (s *Service) GetMockInterviewByID(id int64) (*MockInterview, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	return s.repo.FindByID(id)
}

// ListMockInterviews 分页查询模拟面试列表（管理员）
func (s *Service) ListMockInterviews(req *ListMockInterviewRequest) (*response.PageResponse[MockInterview], error) {
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

	mockInterviews, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}

	records := make([]MockInterview, 0, len(mockInterviews))
	for _, mockInterview := range mockInterviews {
		records = append(records, *mockInterview)
	}

	return &response.PageResponse[MockInterview]{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		Records:  records,
	}, nil
}

// ListMyMockInterviews 分页查询我的模拟面试列表
func (s *Service) ListMyMockInterviews(req *ListMockInterviewRequest, userID int64) (*response.PageResponse[MockInterview], error) {
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return nil, errors.New("无效的用户ID")
	}
	req.UserID = userID
	return s.ListMockInterviews(req)
}
