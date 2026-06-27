package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// Client 是 ES 服务的客户端封装
type Client struct {
	client *elasticsearch.Client // 官方库的客户端
}

type SearchResult struct {
	Total int64
	Hits  []SearchHit
}

type SearchHit struct {
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

// NewClient 创建一个新的 ES 客户端
func NewClient(addresses []string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
	}, nil
}

// IndexDocument 将文档索引到指定的 ES 索引中
// index: ES 索引名称
// id: 文档 ID
// doc: 要索引的文档对象，可以是任意结构体或 map
func (c *Client) IndexDocument(index string, id int64, doc any) error {
	// 1. 把 doc 序列化成 JSON
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("ES 序列化失败：%w", err)
	}
	// 2. 构造索引请求
	req := esapi.IndexRequest{
		Index:      index,                     // "questions" / "posts"
		DocumentID: strconv.FormatInt(id, 10), // int64 → "1"
		Body:       bytes.NewReader(body),     // []byte → io.Reader
		Refresh:    "true",                    // 写入后立即刷新可搜索
	}
	// 3. 发送请求
	res, err := req.Do(context.Background(), c.client)
	if err != nil {
		return fmt.Errorf("ES 索引请求失败:%w", err)
	}
	defer res.Body.Close()
	// 4. 检查响应状态码
	if res.IsError() {
		return fmt.Errorf("ES 索引请求失败: %s", res.Status())
	}
	return nil
}

// Search 在指定的 ES 索引中执行搜索查询
// index: ES 索引名称
// query: 搜索查询的 JSON 字节数组
func (c *Client) Search(index string, query []byte) (*SearchResult, error) {
	// 1. 构造搜索请求
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  bytes.NewReader(query),
	}

	// 2. 发送请求
	res, err := req.Do(context.Background(), c.client)
	if err != nil {
		return nil, fmt.Errorf("ES 搜索请求失败: %w", err)
	}
	defer res.Body.Close()

	// 3. 检查响应状态码
	if res.IsError() {
		return nil, fmt.Errorf("ES 搜索请求失败: %s", res.Status())
	}

	// 4. 解析响应体
	var esResponse struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []SearchHit `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return nil, fmt.Errorf("ES 响应解析失败: %w", err)
	}

	return &SearchResult{
		Total: esResponse.Hits.Total.Value,
		Hits:  esResponse.Hits.Hits,
	}, nil
}

// DeleteDocument 从指定的 ES 索引中删除文档
// index: ES 索引名称
// documentID: 要删除的文档 ID
func (c *Client) DeleteDocument(index string, documentID string) error {
	// 1. 构造删除请求
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    "true", // 删除后立即刷新可搜索
	}

	// 2. 发送请求
	res, err := req.Do(context.Background(), c.client)
	if err != nil {
		return fmt.Errorf("ES 删除请求失败: %w", err)
	}
	defer res.Body.Close()

	// 3. 检查响应状态码
	if res.IsError() {
		return fmt.Errorf("ES 删除请求失败: %s", res.Status())
	}

	return nil
}
