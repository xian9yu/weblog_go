package article

import (
	"weblog/models"

	"gorm.io/gorm"
)

type Controller struct {
	articleRepo *models.ArticleRepository // 只拿自己需要的 repo
	userRepo    *models.UserRepository    // 👍 新增：用户的 Repo
	baseDB      *gorm.DB
}

func NewController(repo *models.ArticleRepository, ur *models.UserRepository, db *gorm.DB) *Controller {
	return &Controller{
		articleRepo: repo,
		userRepo:    ur,
		baseDB:      db,
	}
}
