package mockinterview

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
