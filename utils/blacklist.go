package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// MemoryBlacklist token 黑名单
type MemoryBlacklist struct {
	mu    sync.RWMutex
	store map[string]time.Time // key: token_md5, value: 过期时间
}

// Blacklist 全局单例
var Blacklist = &MemoryBlacklist{
	store: make(map[string]time.Time),
}

// init Go语言原生的初始化函数，在项目启动时自动开启后台清理协程
func init() {
	go Blacklist.GcLoop()
}

// Add 将Token加入黑名单
func (b *MemoryBlacklist) Add(tokenMd5 string, duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.store[tokenMd5] = time.Now().Add(duration)
}

// Contains 检查Token是否在黑名单中（纯内存查询）
func (b *MemoryBlacklist) Contains(tokenMd5 string) bool {
	b.mu.RLock()
	expireAt, exists := b.store[tokenMd5]
	b.mu.RUnlock()

	if !exists {
		return false
	}

	// 如果查询时发现过期了，顺手清理掉
	if time.Now().After(expireAt) {
		b.mu.Lock()
		delete(b.store, tokenMd5)
		b.mu.Unlock()
		return false
	}
	return true
}

// BlockCurrentToken 职责：从请求头自动捞出 Token，计算残余寿命并强给 +1 分钟送进黑名单
func (b *MemoryBlacklist) BlockCurrentToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return
	}
	tokenStr := parts[1]

	// 解密 Token
	claims, err := ParseToken(tokenStr)
	if err != nil {
		return // 已经过期的本来就无法通过中间件，不需要处理
	}

	// 计算剩余寿命
	remainDuration := time.Until(claims.ExpiresAt.Time)
	if remainDuration > 0 {
		hasher := md5.New()
		hasher.Write([]byte(tokenStr))
		tokenMd5 := hex.EncodeToString(hasher.Sum(nil))

		// 强给 +1 分钟防御
		minutes := uint64(remainDuration.Minutes()) + 1

		// 🌟 修正点 1：使用当前绑定的实例指针 b，而不是死死捆绑全局变量 Blacklist
		b.Add(tokenMd5, time.Duration(minutes)*time.Minute)
	}
}

// GcLoop 🌟 行业标准补丁：主动定时清理器（每10分钟主动清洗一次内存）
// 防止那些被拉黑后、再也没发起过请求的Token永久堆积在内存里
func (b *MemoryBlacklist) GcLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		b.mu.Lock()
		now := time.Now()
		for key, expireAt := range b.store {
			if now.After(expireAt) {
				delete(b.store, key)
			}
		}
		b.mu.Unlock()
	}
}
