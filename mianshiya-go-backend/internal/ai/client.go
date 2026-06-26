package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ChatMessage 表示一条对话消息
// Role 取值："system"（系统指令）、"user"（用户）、"assistant"（AI回复）
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 请求体
type chatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// 响应体
type chatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// Client 是 AI 服务的 HTTP 客户端
type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewClient 创建一个新的 AI 客户端
// apiKey: 火山引擎的 API Key
// baseURL: 火山引擎 Ark 的 API 地址，如 https://ark.cn-beijing.volces.com/api/v3
// model: 使用的模型名称，如 deepseek-v3-241226
func NewClient(apiKey, baseURL, model string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // AI 回复可能需要较长时间，设 2 分钟超时
		},
	}
}

// Chat 发送消息列表给 AI，返回 AI 的回复文本
// messages: 完整的对话历史，包含 system/user/assistant 角色
func (c *Client) Chat(messages []ChatMessage) (string, error) {
	// 1. 构造请求体
	reqBody := chatRequest{
		Model:    c.model,
		Messages: messages,
	}

	// 2. 把 Go 结构体序列化为 JSON 字节数组
	// json.Marshal 的作用：Go结构体 → JSON格式的 []byte
	// 例如 chatRequest{Model:"doubao", Messages:[...]}
	//   → []byte(`{"model":"doubao","messages":[...]}`)
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 3. 创建 HTTP POST 请求
	// c.baseURL 是 "https://ark.cn-beijing.volces.com/api/v3"
	// 拼上 "/chat/completions" 就是完整地址
	url := c.baseURL + "/chat/completions"
	// bytes.NewReader(jsonData): 把 []byte 包装成 io.Reader 接口
	// http.NewRequest 的第三个参数 body 需要 io.Reader 类型
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 4. 设置请求头
	// Content-Type 告诉服务端我们发的是 JSON
	req.Header.Set("Content-Type", "application/json")
	// Authorization 是 Bearer Token 认证方式，格式固定为 "Bearer <apiKey>"
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 5. 发送 HTTP 请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求AI服务失败: %w", err)
	}
	// defer: 这行代码不会立即执行，而是在当前函数 Chat() return 之前执行
	// resp.Body.Close() 关闭响应体，释放 HTTP 连接资源
	defer resp.Body.Close()

	// 6. 读取响应体的全部内容
	// io.ReadAll: 从 io.Reader（这里是 resp.Body）读取所有数据到 []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 7. 检查 HTTP 状态码，200 才表示成功
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI服务返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	// 8. 把 JSON 响应反序列化为 Go 结构体
	// json.Unmarshal: JSON格式的 []byte → Go结构体，和 json.Marshal 是相反的操作
	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 9. 检查 AI 是否返回了内容
	if len(chatResp.Choices) == 0 {
		return "", errors.New("AI 调用失败，没有返回结果")
	}

	// 10. 取出第一条回复的文本内容并返回
	return chatResp.Choices[0].Message.Content, nil
}
