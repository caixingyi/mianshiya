package post

import (
	"encoding/json"
	"errors"
	"log"
	"mianshiya-go-backend/internal/es"
	"mianshiya-go-backend/internal/response"
	"mianshiya-go-backend/internal/user"
	"strconv"

	"github.com/sony/gobreaker"
)

// Service 帖子服务层
type Service struct {
	repo      *Repository
	userSvc   *user.Service
	es        *es.Client
	esBreaker *gobreaker.CircuitBreaker
}

// NewService 创建帖子服务实例
func NewService(r *Repository, userSvc *user.Service, esClient *es.Client, esBreaker *gobreaker.CircuitBreaker) *Service {
	return &Service{repo: r, userSvc: userSvc, es: esClient, esBreaker: esBreaker}
}

// 转换 Post 到 PostResponse
func (s *Service) toPostResponse(post *Post, loginUserID int64) (*PostResponse, error) {
	tagList := make([]string, 0)
	if post.Tags != "" {
		if err := json.Unmarshal([]byte(post.Tags), &tagList); err != nil {
			return nil, err
		}
	}
	hasThumb := false
	hasFavour := false
	if loginUserID > 0 {
		if ht, err := s.repo.HasThumb(post.ID, loginUserID); err == nil {
			hasThumb = ht
		}
		if hf, err := s.repo.HasFavour(post.ID, loginUserID); err == nil {
			hasFavour = hf
		}
	}
	return &PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		TagList:   tagList,
		ThumbNum:  post.ThumbNum,
		FavourNum: post.FavourNum,
		HasThumb:  hasThumb,
		HasFavour: hasFavour,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

// AddPost 添加帖子
func (s *Service) AddPost(req *AddPostRequest, userID int64) (int64, error) {
	// 参数校验
	if req == nil {
		return 0, errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return 0, errors.New("无效的用户ID")
	}
	if req.Title == "" {
		return 0, errors.New("帖子标题不能为空")
	}
	if req.Content == "" {
		return 0, errors.New("帖子内容不能为空")
	}
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}
	tagsBytes, err := json.Marshal(tags)
	if err != nil {
		return 0, errors.New("标签格式错误")
	}

	// 创建 Post 实例
	post := &Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
		Tags:    string(tagsBytes),
	}

	// 调用 Repository 添加帖子
	postID, err := s.repo.Create(post)
	if err != nil {
		return 0, err
	}

	// 异步添加到 Elasticsearch
	go func() {
		if err := s.es.IndexDocument("posts", postID, post); err != nil {
			// 记录日志或处理错误
			log.Printf("Failed to index post in Elasticsearch: %v", err)
		}
	}()

	return postID, nil
}

// DeletePost 删除帖子
func (s *Service) DeletePost(req *DeletePostRequest, userID int64) error {
	// 参数校验
	if req == nil {
		return errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return errors.New("无效的用户ID")
	}
	if req.ID <= 0 {
		return errors.New("无效的帖子ID")
	}
	// 验证帖子存在且用户有权限删除
	post, err := s.repo.FindByID(req.ID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	isAdmin, err := s.userSvc.IsAdmin(userID)
	if err != nil {
		return errors.New("无法验证用户权限")
	}
	if post.UserID != userID && !isAdmin {
		return errors.New("无权限删除该帖子")
	}

	// 调用 Repository 删除帖子
	err = s.repo.DeleteByID(req.ID)
	if err != nil {
		return errors.New("删除帖子失败")
	}

	// 异步删除 Elasticsearch 中的帖子
	go func() {
		if err := s.es.DeleteDocument("posts", strconv.FormatInt(req.ID, 10)); err != nil {
			// 记录日志或处理错误
			log.Printf("Failed to delete post from Elasticsearch: %v", err)
		}
	}()
	return nil
}

// GetPostByID 获取帖子详情，loginUserID 为 0 表示未登录（不查点赞/收藏状态）
func (s *Service) GetPostByID(id int64, loginUserID int64) (*PostResponse, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp, err := s.toPostResponse(post, loginUserID)
	if err != nil {
		return nil, err
	}
	// 关联查发帖人
	u, _ := s.userSvc.GetUserResponseByID(post.UserID)
	if u != nil {
		resp.User = u
	}
	return resp, nil
}

// ListPosts 分页查询帖子列表（返回 VO），loginUserID 为 0 表示未登录
func (s *Service) ListPosts(req *ListPostsRequest, loginUserID int64) (*response.PageResponse[PostResponse], error) {
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
	posts, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}
	responses := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		resp, err := s.toPostResponse(post, loginUserID)
		if err != nil {
			return nil, err
		}
		u, _ := s.userSvc.GetUserResponseByID(post.UserID)
		if u != nil {
			resp.User = u
		}
		responses = append(responses, *resp)
	}
	return &response.PageResponse[PostResponse]{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		Records:  responses,
	}, nil
}

// ListPostPage 管理员分页查询帖子列表（返回实体）
func (s *Service) ListPostPage(req *ListPostsRequest) (*response.PageResponse[Post], error) {
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
	posts, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}
	records := make([]Post, 0, len(posts))
	for _, post := range posts {
		records = append(records, *post)
	}
	return &response.PageResponse[Post]{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		Records:  records,
	}, nil
}

// ListMyPosts 获取我的帖子列表
func (s *Service) ListMyPostsVO(req *ListPostsRequest, userID int64) (*response.PageResponse[PostResponse], error) {
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return nil, errors.New("无效的用户ID")
	}
	req.UserID = userID
	return s.ListPosts(req, userID)
}

// UpdatePost 更新帖子（管理员接口）
func (s *Service) UpdatePost(req *UpdatePostRequest, userID int64) error {
	// 参数校验
	if req == nil {
		return errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return errors.New("无效的用户ID")
	}
	if req.ID <= 0 {
		return errors.New("无效的帖子ID")
	}

	// 构造更新数据
	updates := make(map[string]any)
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return errors.New("标签格式错误")
		}
		updates["tags"] = string(tagsBytes)
	}
	// 调用 Repository 更新帖子
	err := s.repo.UpdateByID(req.ID, updates)
	if err != nil {
		return errors.New("更新帖子失败")
	}

	// 异步更新 Elasticsearch 中的帖子
	go func() {
		post, err := s.repo.FindByID(req.ID)
		if err != nil {
			log.Printf("Failed to find post: %v", err)
			return
		}
		if err := s.es.IndexDocument("posts", req.ID, post); err != nil {
			log.Printf("Failed to index post in Elasticsearch: %v", err)
		}
	}()

	return nil
}

// EditPost 编辑帖子
func (s *Service) EditPost(req *EditPostRequest, userID int64) error {
	// 参数校验
	if req == nil {
		return errors.New("请求参数不能为空")
	}
	if userID <= 0 {
		return errors.New("无效的用户ID")
	}
	if req.ID <= 0 {
		return errors.New("无效的帖子ID")
	}
	// 验证帖子存在且用户有权限更新
	post, err := s.repo.FindByID(req.ID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	isAdmin, err := s.userSvc.IsAdmin(userID)
	if err != nil {
		return errors.New("无法验证用户权限")
	}
	if post.UserID != userID && !isAdmin {
		return errors.New("无权限编辑该帖子")
	}
	// 构造更新数据
	updates := make(map[string]any)
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return errors.New("标签格式错误")
		}
		updates["tags"] = string(tagsBytes)
	}

	// 调用 Repository 更新帖子
	err = s.repo.UpdateByID(post.ID, updates)
	if err != nil {
		return errors.New("更新帖子失败")
	}

	// 异步更新 Elasticsearch 中的帖子
	go func() {
		updatedPost, err := s.repo.FindByID(post.ID)
		if err != nil {
			log.Printf("查询帖子失败: %v", err)
			return
		}
		if err := s.es.IndexDocument("posts", post.ID, updatedPost); err != nil {
			log.Printf("索引帖子到 Elasticsearch 失败: %v", err)
		}
	}()
	return nil
}

// searchFromMySQL MySQL LIKE 搜索，作为 ES 的降级方案
func (s *Service) searchFromMySQL(keyword string, current, pageSize int64) (*response.PageResponse[PostResponse], error) {
	req := &ListPostsRequest{
		Current:    current,
		PageSize:   pageSize,
		SearchText: keyword,
	}
	return s.ListPosts(req, 0) // loginUserID=0 表示未登录，不查点赞/收藏状态
}

func (s *Service) SearchPosts(keyword string, current, pageSize int64) (*response.PageResponse[PostResponse], error) {
	// 1. 校验参数
	if current <= 0 {
		current = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 200 {
		return nil, errors.New("参数错误")
	}

	// 空关键词不查 ES，直接走普通列表
	if keyword == "" {
		return s.ListPosts(&ListPostsRequest{
			Current:  current,
			PageSize: pageSize,
		}, 0)
	}

	// 2. 构造 ES 查询 DSL（JSON）
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"title", "content", "tags"},
			},
		},
		"from": (current - 1) * pageSize,
		"size": pageSize,
	}
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, errors.New("构造 Elasticsearch 查询失败")
	}
	// 先尝试从 Elasticsearch 搜索
	resultAny, err := s.esBreaker.Execute(func() (any, error) {
		return s.es.Search("posts", queryBytes)
	})
	if err != nil {
		log.Printf("Elasticsearch 搜索失败或熔断打开: %v", err)
		// 如果 Elasticsearch 搜索失败，降级到 MySQL 搜索
		return s.searchFromMySQL(keyword, current, pageSize)
	}
	esResults := resultAny.(*es.SearchResult)
	// 3. 解析 ES 搜索结果
	records := make([]PostResponse, 0, len(esResults.Hits))
	for _, hit := range esResults.Hits {
		var post Post
		if err := json.Unmarshal(hit.Source, &post); err != nil {
			log.Printf("解析 Elasticsearch 搜索结果失败: %v", err)
			continue
		}
		resp, err := s.toPostResponse(&post, 0) // loginUserID=0 表示未登录，不查点赞/收藏状态
		if err != nil {
			log.Printf("解析帖子响应失败: %v", err)
			continue
		}
		u, _ := s.userSvc.GetUserResponseByID(post.UserID)
		if u != nil {
			resp.User = u
		}
		records = append(records, *resp)
	}
	return &response.PageResponse[PostResponse]{
		Current:  current,
		PageSize: pageSize,
		Total:    esResults.Total,
		Records:  records,
	}, nil
}
