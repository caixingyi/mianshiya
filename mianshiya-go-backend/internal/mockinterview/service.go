package mockinterview

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"mianshiya-go-backend/internal/ai"
	"mianshiya-go-backend/internal/response"
	"mianshiya-go-backend/internal/user"
)

// Service 模拟面试服务层
type Service struct {
	repo    *Repository
	userSvc *user.Service
	ai      *ai.Client // 新增：AI 客户端
}

// NewService 创建模拟面试服务实例
func NewService(repo *Repository, userSvc *user.Service, aiClient *ai.Client) *Service {
	return &Service{repo: repo, userSvc: userSvc, ai: aiClient}
}

// ======================== 基础 CRUD（已有） ========================

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

// ======================== AI 对话事件处理 ========================

// HandleEvent 处理模拟面试事件（start/chat/end），返回 AI 的回复文本
func (s *Service) HandleEvent(req *HandleEventRequest, userID int64) (string, error) {
	if req == nil {
		return "", errors.New("请求参数不能为空")
	}
	if req.ID <= 0 {
		return "", errors.New("参数错误")
	}
	if userID <= 0 {
		return "", errors.New("无效的用户ID")
	}

	// 1. 查询面试记录
	mockInterview, err := s.repo.FindByID(req.ID)
	if err != nil {
		return "", fmt.Errorf("模拟面试不存在: %w", err)
	}

	// 2. 校验所有权：只能操作自己的面试
	if mockInterview.UserID != userID {
		return "", errors.New("无权限操作该模拟面试")
	}

	// 3. 根据事件类型分发
	switch req.Event {
	case EventStart:
		return s.handleStartEvent(mockInterview)
	case EventChat:
		return s.handleChatEvent(req, mockInterview)
	case EventEnd:
		return s.handleEndEvent(mockInterview)
	default:
		return "", fmt.Errorf("不支持的事件类型: %s", req.Event)
	}
}

// handleStartEvent 处理"开始"事件：构建 system prompt → 调用 AI → 更新状态为进行中
func (s *Service) handleStartEvent(mockInterview *MockInterview) (string, error) {
	// 1. 构建 system prompt
	// 把用户的工作经验、岗位、难度填入提示词，让 AI 扮演面试官
	systemPrompt := fmt.Sprintf(
		"你是一位严厉的程序员面试官，我是候选人，来应聘 %s 的 %s 岗位，面试难度为 %s。"+
			"请你向我依次提出问题（最多 20 个问题），我也会依次回复。在这期间请完全保持真人面试官的口吻，"+
			"比如适当引导学员、或者表达出你对学员回答的态度。\n"+
			"必须满足如下要求：\n"+
			"1. 当学员回复 \"开始\" 时，你要正式开始面试\n"+
			"2. 当学员表示希望 \"结束面试\" 时，你要结束面试\n"+
			"3. 此外，当你觉得这场面试可以结束时"+
			"（比如候选人回答结果较差、不满足工作年限的招聘需求、或者候选人态度不礼貌），"+
			"必须主动提出面试结束，不用继续询问更多问题了。并且要在回复中包含字符串【面试结束】\n"+
			"4. 面试结束后，应该给出候选人整场面试的表现和总结。",
		mockInterview.WorkExperience,
		mockInterview.JobPosition,
		mockInterview.Difficulty,
	)

	// 2. 构造初始消息列表：system指令 + 用户说"开始"
	aiMessages := []ai.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "开始"},
	}

	// 3. 调用 AI
	answer, err := s.ai.Chat(aiMessages)
	if err != nil {
		return "", fmt.Errorf("AI 调用失败: %w", err)
	}

	// 4. 把 AI 的回复追加到消息列表
	aiMessages = append(aiMessages, ai.ChatMessage{Role: "assistant", Content: answer})

	// 5. 转换为 DB 存储格式（role/content → role/message），序列化为 JSON
	messagesJSON := convertToChatMessages(aiMessages)
	jsonStr, err := json.Marshal(messagesJSON)
	if err != nil {
		return "", fmt.Errorf("序列化消息失败: %w", err)
	}

	// 6. 更新数据库：消息记录 + 状态改为"进行中"
	_ = s.repo.UpdateByID(mockInterview.ID, map[string]any{
		"messages": string(jsonStr),
		"status":   StatusInProgress,
	})

	return answer, nil
}

// handleChatEvent 处理"对话"事件：加载历史 → 追加用户消息 → 调用 AI → 保存
func (s *Service) handleChatEvent(req *HandleEventRequest, mockInterview *MockInterview) (string, error) {
	// 1. 从 DB 加载历史消息，转换为 AI 格式
	aiMessages, err := loadHistoryMessages(mockInterview.Messages)
	if err != nil {
		return "", fmt.Errorf("加载历史消息失败: %w", err)
	}

	// 2. 追加当前用户消息
	aiMessages = append(aiMessages, ai.ChatMessage{Role: "user", Content: req.Message})

	// 3. 调用 AI
	answer, err := s.ai.Chat(aiMessages)
	if err != nil {
		return "", fmt.Errorf("AI 调用失败: %w", err)
	}

	// 4. 追加 AI 回复
	aiMessages = append(aiMessages, ai.ChatMessage{Role: "assistant", Content: answer})

	// 5. 序列化消息并保存
	messagesJSON := convertToChatMessages(aiMessages)
	jsonStr, err := json.Marshal(messagesJSON)
	if err != nil {
		return "", fmt.Errorf("序列化消息失败: %w", err)
	}

	updates := map[string]any{
		"messages": string(jsonStr),
	}
	// 如果 AI 回复中包含"【面试结束】"，自动结束面试
	if containsEndMarker(answer) {
		updates["status"] = StatusEnded
	}
	_ = s.repo.UpdateByID(mockInterview.ID, updates)

	return answer, nil
}

// handleEndEvent 处理"结束"事件：加载历史 → 追加"结束" → 调用 AI 总结 → 更新状态
func (s *Service) handleEndEvent(mockInterview *MockInterview) (string, error) {
	// 1. 从 DB 加载历史消息
	aiMessages, err := loadHistoryMessages(mockInterview.Messages)
	if err != nil {
		return "", fmt.Errorf("加载历史消息失败: %w", err)
	}

	// 2. 追加结束消息
	aiMessages = append(aiMessages, ai.ChatMessage{Role: "user", Content: "结束"})

	// 3. 调用 AI
	answer, err := s.ai.Chat(aiMessages)
	if err != nil {
		return "", fmt.Errorf("AI 调用失败: %w", err)
	}

	// 4. 追加 AI 回复
	aiMessages = append(aiMessages, ai.ChatMessage{Role: "assistant", Content: answer})

	// 5. 序列化并保存，状态改为"已结束"
	messagesJSON := convertToChatMessages(aiMessages)
	jsonStr, err := json.Marshal(messagesJSON)
	if err != nil {
		return "", fmt.Errorf("序列化消息失败: %w", err)
	}

	_ = s.repo.UpdateByID(mockInterview.ID, map[string]any{
		"messages": string(jsonStr),
		"status":   StatusEnded,
	})

	return answer, nil
}

// ======================== 辅助函数 ========================

// loadHistoryMessages 把 DB 中的 JSON 消息记录解析为 AI 格式（role/content）
// DB 存储格式：[{"role":"user","message":"你好"}, ...]
// AI API 格式：[{Role:"user", Content:"你好"}, ...]
func loadHistoryMessages(messagesJSON string) ([]ai.ChatMessage, error) {
	// 空消息记录（首次对话）返回空列表
	if messagesJSON == "" || messagesJSON == "[]" {
		return make([]ai.ChatMessage, 0), nil
	}

	// 反序列化 DB 中的 JSON 数组
	var chatMessages []ChatMessage
	if err := json.Unmarshal([]byte(messagesJSON), &chatMessages); err != nil {
		return nil, fmt.Errorf("解析历史消息失败: %w", err)
	}

	// 转换字段名：message → content
	result := make([]ai.ChatMessage, 0, len(chatMessages))
	for _, cm := range chatMessages {
		result = append(result, ai.ChatMessage{
			Role:    cm.Role,
			Content: cm.Message, // ← 字段名转换
		})
	}
	return result, nil
}

// convertToChatMessages 把 AI 格式消息（role/content）转为 DB 存储格式（role/message）
func convertToChatMessages(aiMessages []ai.ChatMessage) []ChatMessage {
	result := make([]ChatMessage, 0, len(aiMessages))
	for _, m := range aiMessages {
		result = append(result, ChatMessage{
			Role:    m.Role,
			Message: m.Content, // ← 字段名转换
		})
	}
	return result
}

// containsEndMarker 检查 AI 回复中是否包含面试结束标记
// Java 项目用 chatAnswer.contains("【面试结束】")
func containsEndMarker(answer string) bool {
	return strings.Contains(answer, "【面试结束】") || strings.Contains(answer, "[面试结束]")
}
