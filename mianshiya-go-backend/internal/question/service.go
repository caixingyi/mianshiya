package question

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"mianshiya-go-backend/internal/ai"
	"mianshiya-go-backend/internal/es"
	"mianshiya-go-backend/internal/response"

	"github.com/sony/gobreaker"
)

// Service 题目服务层
type Service struct {
	repo      *Repository
	ai        *ai.Client                // AI 客户端，用于生成题目和题解
	es        *es.Client                // ES 客户端，用于搜索题目
	esBreaker *gobreaker.CircuitBreaker // ES 熔断器，用于保护 ES 服务
}

// NewService 创建题目服务实例
func NewService(r *Repository, aiClient *ai.Client, esClient *es.Client, esBreaker *gobreaker.CircuitBreaker) *Service {
	return &Service{repo: r, ai: aiClient, es: esClient, esBreaker: esBreaker}
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
// req: 添加题目的请求参数
// userID: 当前登录用户的 ID
// 返回值: 新增题目的 ID，或错误信息
// 该方法会将题目保存到数据库，并异步同步到 Elasticsearch
// 如果 Elasticsearch 同步失败，不会影响主流程
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
	id, err := s.repo.Create(question)
	if err != nil {
		return 0, err
	}
	// 双写 ES(异步，失败不影响主流程)
	question.ID = id // 确保 ID 已设置
	go func() {
		if err := s.es.IndexDocument("questions", id, question); err != nil {
			log.Printf("[ES] 同步题目到 Elasticsearch 失败: %v", err)
		}
	}()

	return id, nil
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

	// 从MySQL中删除题目
	err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}

	// 异步删除 Elasticsearch 中的题目
	go func() {
		if err := s.es.DeleteDocument("questions", strconv.FormatInt(id, 10)); err != nil {
			log.Printf("[ES] 删除题目从 Elasticsearch 失败: %v", err)
		}
	}()
	return nil
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
	// 更新Mysql数据库
	err := s.repo.UpdateByID(req.ID, updates)
	if err != nil {
		return err
	}
	// 异步更新 Elasticsearch
	go func() {
		question, err := s.repo.FindByID(req.ID)
		if err != nil {
			log.Printf("[ES] 更新题目到 Elasticsearch 失败: %v", err)
			return
		}
		if err := s.es.IndexDocument("questions", req.ID, question); err != nil {
			log.Printf("[ES] 更新题目到 Elasticsearch 失败: %v", err)
		}
	}()
	return nil
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
	if err := s.repo.DeleteBatchByIDs(req.QuestionIDList); err != nil {
		return err
	}

	// 双写 ES：逐条异步删除
	for _, id := range req.QuestionIDList {
		go func(questionID int64) {
			if err := s.es.DeleteDocument("questions", strconv.FormatInt(questionID, 10)); err != nil {
				log.Printf("[ES] 删除题目失败 [%d]: %v", questionID, err)
			}
		}(id)
	}

	return nil
}

// SearchQuestions 用 ES 全文搜索题目，ES 失败降级到 MySQL
// keyword: 搜索关键词
// current/pageSize: 分页
func (s *Service) SearchQuestions(keyword string, current, pageSize int64) (*response.PageResponse[QuestionResponse], error) {
	if current <= 0 {
		current = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 200 {
		return nil, errors.New("参数错误")
	}

	// 1. 构造 ES 查询 DSL（JSON）
	// multi_match 在 title 和 content 两个字段里搜 keyword
	query := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  keyword,
				"fields": []string{"title", "content", "answer"},
			},
		},
		"from": (current - 1) * pageSize, // 偏移量
		"size": pageSize,                 // 返回数量
	}
	queryJSON, _ := json.Marshal(query)

	// 2. 调 ES 搜索
	resultAny, err := s.esBreaker.Execute(func() (any, error) {
		return s.es.Search("questions", queryJSON)
	})
	if err != nil {
		// ES 挂了，降级到 MySQL LIKE
		log.Printf("[ES] 搜索失败或熔断打开，降级到 MySQL: %v", err)
		return s.searchFromMySQL(keyword, current, pageSize)
	}
	result := resultAny.(*es.SearchResult)
	// 3. 解析 ES 命中的文档
	records := make([]QuestionResponse, 0, len(result.Hits))
	for _, hit := range result.Hits {
		var q Question
		if err := json.Unmarshal(hit.Source, &q); err != nil {
			continue
		}
		resp, _ := s.toQuestionResponse(&q)
		if resp != nil {
			records = append(records, *resp)
		}
	}

	return &response.PageResponse[QuestionResponse]{
		Current:  current,
		PageSize: pageSize,
		Total:    result.Total,
		Records:  records,
	}, nil
}

// searchFromMySQL MySQL LIKE 搜索，作为 ES 的降级方案
func (s *Service) searchFromMySQL(keyword string, current, pageSize int64) (*response.PageResponse[QuestionResponse], error) {
	req := &ListQuestionRequest{
		Current:    current,
		PageSize:   pageSize,
		SearchText: keyword,
	}
	return s.ListQuestions(req)
}

// ======================== AI 生成题目 ========================

// AIGenerateQuestions AI 生成题目，对应 Java 的 aiGenerateQuestions
// questionType: 题目方向，如 "Java"
// number: 生成数量，如 10
// userID: 创建者 ID
func (s *Service) AIGenerateQuestions(questionType string, number int, userID int64) error {
	// 1. 参数校验
	if questionType == "" {
		return errors.New("题目类型不能为空")
	}
	if number <= 0 {
		number = 10 // 默认 10 道
	}
	if userID <= 0 {
		return errors.New("无效的用户ID")
	}

	// 2. 构建 system prompt — 告诉 AI 它是什么角色、输出什么格式
	systemPrompt := fmt.Sprintf(
		"你是一位专业的程序员面试官，你要帮我生成 %d 道 %s 面试题，要求输出格式如下：\n\n"+
			"1. 什么是 Java 中的反射？\n"+
			"2. Java 8 中的 Stream API 有什么作用？\n"+
			"3. xxxxxx\n\n"+
			"除此之外，请不要输出任何多余的内容，不要输出开头、也不要输出结尾，只输出上面的列表。\n\n"+
			"接下来我会给你要生成的题目数量和题目方向",
		number, questionType,
	)

	// 3. 构建 user prompt
	userPrompt := fmt.Sprintf("题目数量：%d, 题目方向：%s", number, questionType)

	// 4. 调用 AI 生成题目列表
	answer, err := s.ai.Chat([]ai.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	})
	if err != nil {
		return fmt.Errorf("AI 生成题目失败: %w", err)
	}

	// 5. 解析 AI 返回的题目列表
	// AI 返回格式：每行一个题目，如 "1. 什么是反射？"
	lines := strings.Split(answer, "\n")
	titles := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue // 跳过空行
		}
		// 去掉序号前缀 "1. "、"2. " 等
		// 找到第一个空格或 ". " 后的内容
		title := removeNumberPrefix(line)
		// 去掉反引号（AI 偶尔会输出 markdown 代码标记）
		title = strings.ReplaceAll(title, "`", "")
		title = strings.TrimSpace(title)
		if title != "" {
			titles = append(titles, title)
		}
	}

	if len(titles) == 0 {
		return errors.New("AI 未生成任何题目")
	}

	// 6. 逐题生成题解（带重试），构建 Question 列表
	questions := make([]*Question, 0, len(titles))
	for _, title := range titles {
		answer, err := s.generateAnswerWithRetry(title, 3)
		if err != nil {
			log.Printf("[AI题解] 重试3次后仍失败 [%s]: %v", title, err)
			continue // 跳过这道题，不入库
		}
		log.Printf("[AI题解] 成功 [%s]: answer长度=%d", title, len(answer))
		questions = append(questions, &Question{
			Title:  title,
			Answer: answer,
			Tags:   `["待审核"]`, // 和 Java 一致，打上"待审核"标签
			UserID: userID,
		})
	}

	// 7. 批量入库
	if err := s.repo.BatchCreate(questions); err != nil {
		return err
	}

	// 8. 双写 ES（异步逐条同步）
	for _, q := range questions {
		go func(question *Question) {
			if err := s.es.IndexDocument("questions", question.ID, question); err != nil {
				log.Printf("[ES] 同步题目失败 [%d]: %v", question.ID, err)
			}
		}(q)
	}

	return nil

}

// generateAnswer 为题目生成题解，对应 Java 的 aiGenerateQuestionAnswer
func (s *Service) generateAnswer(title string) (string, error) {
	systemPrompt := "你是一位专业的程序员面试官，我会给你一道面试题，请帮我生成详细的题解。要求如下：\n\n" +
		"1. 题解的语句要自然流畅\n" +
		"2. 题解可以先给出总结性的回答，再详细解释\n" +
		"3. 要使用 Markdown 语法输出\n\n" +
		"除此之外，请不要输出任何多余的内容，不要输出开头、也不要输出结尾，只输出题解。\n\n" +
		"接下来我会给你要生成的面试题"

	userPrompt := fmt.Sprintf("面试题：%s", title)

	return s.ai.Chat([]ai.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	})
}

// generateAnswerWithRetry 带重试的题解生成
// maxRetries: 最大重试次数（不含首次调用）
func (s *Service) generateAnswerWithRetry(title string, maxRetries int) (string, error) {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			// 重试前等待，递增延迟：1s → 2s → 3s
			time.Sleep(time.Duration(i) * time.Second)
			log.Printf("[AI题解] 第%d次重试 [%s]", i, title)
		}

		answer, err := s.generateAnswer(title)
		if err != nil {
			lastErr = err
			continue
		}
		if answer == "" {
			lastErr = errors.New("AI 返回空内容")
			continue
		}
		return answer, nil
	}
	return "", fmt.Errorf("重试%d次后仍失败: %w", maxRetries, lastErr)
}

// removeNumberPrefix 去掉题目行前面的序号，如 "1. "、"1、"、"1."
// 对应 Java 的 StrUtil.removePrefix(line, StrUtil.subBefore(line, " ", false))
func removeNumberPrefix(line string) string {
	// 先转换成 rune 切片，方便处理中文字符
	runes := []rune(line)
	i := 0
	// 跳过数字
	for i < len(runes) && runes[i] >= '0' && runes[i] <= '9' {
		i++
	}
	// 跳过分隔符：. 、 . 、 （中英文句号、顿号）
	if i < len(runes) && (runes[i] == '.' || runes[i] == '、' || runes[i] == '。') {
		i++
	}
	// 跳过空格
	if i < len(runes) && runes[i] == ' ' {
		i++
	}
	if i >= len(runes) {
		return line
	}
	return string(runes[i:])
}
