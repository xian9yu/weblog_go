package models

import (
	"gorm.io/gorm"
)

type Repositories struct {
	Article      *ArticleRepository
	Auth         *AuthRepository
	User         *UserRepository
	SystemConfig *SystemConfigRepository
}

// NewRepositories 一次性把所有 Repo 初始化好
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Article:      NewArticleRepository(db),
		Auth:         NewAuthRepository(db),
		User:         NewUserRepository(db),
		SystemConfig: NewSystemConfigRepository(db),
	}
}

// WithTx 集体切换到事务状态
func (repos *Repositories) WithTx(tx *gorm.DB) *Repositories {
	return &Repositories{
		Article:      repos.Article.WithTx(tx),
		Auth:         repos.Auth.WithTx(tx),
		SystemConfig: repos.SystemConfig.WithTx(tx),
		User:         repos.User.WithTx(tx),
	}
}
