package article

import (
	"weblog/models"

	"gorm.io/gorm"
)

//	type Controller struct {
//		articleRepo *models.ArticleRepository // 只拿自己需要的 repo
//		baseDB      *gorm.DB
//		userRepo    *models.UserRepository // 👍 新增：用户的 Repo
//	}
//
//	func NewController(repo *models.ArticleRepository, ur *models.UserRepository, db *gorm.DB) *Controller {
//		return &Controller{
//			articleRepo: repo,
//			baseDB:      db,
//			userRepo:    ur,
//		}
//	}
type Controller struct {
	baseDB *gorm.DB
	repos  *models.Repositories
}

func NewController(repos *models.Repositories) *Controller {
	return &Controller{
		baseDB: repos.Database, // 或者是原生的 db
		repos:  repos,
	}
}
