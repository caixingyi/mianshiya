package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type MemoryTokenStore struct {
	mu            sync.RWMutex
	tokenToUserID map[string]int64
}

func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{
		tokenToUserID: make(map[string]int64),
	}
}

func (s *MemoryTokenStore) CreateToken(ctx context.Context, userID int64) (string, error) {
	// 生成随机 token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	// 存储 token 和 userID 的映射关系
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokenToUserID[token] = userID
	return token, nil
}

func (s *MemoryTokenStore) GetUserID(ctx context.Context, token string) (int64, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, exists := s.tokenToUserID[token]
	return userID, exists, nil
}

func (s *MemoryTokenStore) DeleteToken(ctx context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokenToUserID, token)
	return nil
}

// TokenStore 定义了一个接口，抽象了 token 的创建、获取和删除操作，方便后续替换为 Redis 等持久化存储方案
type TokenStore interface {
	CreateToken(ctx context.Context, userID int64) (string, error)
	GetUserID(ctx context.Context, token string) (int64, bool, error)
	DeleteToken(ctx context.Context, token string) error
}
