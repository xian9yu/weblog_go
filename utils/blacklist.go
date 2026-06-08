package utils

import (
	"sync"
	"time"
)

// MemoryBlacklist token 黑名单
type MemoryBlacklist struct {
	mu    sync.RWMutex
	store map[string]time.Time // key: token_md5, value: 过期时间
}

var Blacklist = &MemoryBlacklist{
	store: make(map[string]time.Time),
}

// Add 将Token加入黑名单
func (b *MemoryBlacklist) Add(tokenMd5 string, duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.store[tokenMd5] = time.Now().Add(duration)
}

// Contains 检查Token是否在黑名单中（纯内存查询，纳秒级响应）
func (b *MemoryBlacklist) Contains(tokenMd5 string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	expireAt, exists := b.store[tokenMd5]
	if !exists {
		return false
	}

	// 如果内存里的黑名单时间已经过了，顺手清理掉
	if time.Now().After(expireAt) {
		b.mu.Lock()
		delete(b.store, tokenMd5)
		b.mu.Unlock()
		return false
	}
	return true
}
