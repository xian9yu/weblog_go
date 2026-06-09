package user

import (
	"weblog/models"

	"gorm.io/gorm"
)

type Controller struct {
	articleRepo *models.ArticleRepository // 只拿自己需要的 repo
	baseDB      *gorm.DB
	userRepo    *models.UserRepository // 👍 新增：用户的 Repo
}

func NewController(repo *models.ArticleRepository, ur *models.UserRepository, db *gorm.DB) *Controller {
	return &Controller{
		articleRepo: repo,
		baseDB:      db,
		userRepo:    ur,
	}
}
