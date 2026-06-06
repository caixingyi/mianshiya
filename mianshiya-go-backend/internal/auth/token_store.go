package auth

import (
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

func (s *MemoryTokenStore) CreateToken(userID int64) (string, error) {
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

func (s *MemoryTokenStore) GetUserID(token string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, exists := s.tokenToUserID[token]
	return userID, exists
}

func (s *MemoryTokenStore) DeleteToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokenToUserID, token)
}
