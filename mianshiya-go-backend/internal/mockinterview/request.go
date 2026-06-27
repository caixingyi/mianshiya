package mockinterview

// 事件类型常量，对应 Java 的 MockInterviewEventEnum
const (
	EventStart = "start" // 开始面试
	EventChat  = "chat"  // 对话消息
	EventEnd   = "end"   // 结束面试
)

// ChatMessage 对话消息记录，存在 DB messages 字段里（JSON 数组）
// 注意：这里字段名是 "message"，和 AI API 用的 "content" 不同，service 层负责转换
type ChatMessage struct {
	Role    string `json:"role"`    // system / user / assistant
	Message string `json:"message"` // 消息文本内容
}

// HandleEventRequest 处理模拟面试事件请求，对应 Java 的 MockInterviewEventRequest
type HandleEventRequest struct {
	Event   string `json:"event" binding:"required"` // start / chat / end
	Message string `json:"message"`                  // 用户消息内容，仅 chat 事件需要
	ID      int64  `json:"id" binding:"required"`    // 模拟面试房间 ID
}

// AddMockInterviewRequest 创建模拟面试请求参数
type AddMockInterviewRequest struct {
	WorkExperience string `json:"workExperience" binding:"required"`
	JobPosition    string `json:"jobPosition" binding:"required"`
	Difficulty     string `json:"difficulty" binding:"required"`
}

// DeleteMockInterviewRequest 删除模拟面试请求参数
type DeleteMockInterviewRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// GetMockInterviewRequest 获取模拟面试详情请求参数
type GetMockInterviewRequest struct {
	ID int64 `form:"id" binding:"required"`
}

// ListMockInterviewRequest 分页查询模拟面试请求参数
type ListMockInterviewRequest struct {
	Current        int64  `json:"current"`
	PageSize       int64  `json:"pageSize"`
	ID             int64  `json:"id"`
	WorkExperience string `json:"workExperience"`
	JobPosition    string `json:"jobPosition"`
	Difficulty     string `json:"difficulty"`
	Status         *int   `json:"status"`
	UserID         int64  `json:"userId"`
}
