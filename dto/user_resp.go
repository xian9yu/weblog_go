package dto

import "time"

// UserResponse 安全脱敏后的用户信息响应（用于返回给前端，绝对不带密码字段）
type UserResponse struct {
	ID        uint64    `json:"id" gorm:"column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Email     string    `json:"email" gorm:"column:email"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// LoginResponse 登录成功后吐给前端的黄金组合
type LoginResponse struct {
	Token string       `json:"token"` // JWT 令牌
	User  UserResponse `json:"user"`  // 顺便把用户信息带回去供前端展示
}
