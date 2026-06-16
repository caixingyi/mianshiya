package post

import (
	"encoding/json"
	"errors"
	"mianshiya-go-backend/internal/response"
	"mianshiya-go-backend/internal/user"
)

// Service 帖子服务层
type Service struct {
	repo    *Repository
	userSvc *user.Service
}

// NewService 创建帖子服务实例
func NewService(r *Repository, userSvc *user.Service) *Service {
	return &Service{repo: r, userSvc: userSvc}
}

// 转换 Post 到 PostResponse
func (s *Service) toPostResponse(post *Post) (*PostResponse, error) {
	tagList := make([]string, 0)
	if post.Tags != "" {
		if err := json.Unmarshal([]byte(post.Tags), &tagList); err != nil {
			return nil, err
		}
	}
	return &PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		TagList:   tagList,
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
	return s.repo.DeleteByID(req.ID)
}

// GetPostByID 获取帖子详情
func (s *Service) GetPostByID(id int64) (*PostResponse, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp, err := s.toPostResponse(post)
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

// ListPosts 分页查询帖子列表（返回 VO）
func (s *Service) ListPosts(req *ListPostsRequest) (*response.PageResponse[PostResponse], error) {
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
		resp, err := s.toPostResponse(post)
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
	return s.ListPosts(req)
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
	return s.repo.UpdateByID(req.ID, updates)
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
	return s.repo.UpdateByID(post.ID, updates)
}
