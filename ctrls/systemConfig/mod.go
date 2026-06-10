package systemConfig

import (
	"weblog/models"

	"gorm.io/gorm"
)

//	type Controller struct {
//		systemConfigRepo *models.SystemConfigRepository // 只拿自己需要的 repo
//		baseDB           *gorm.DB
//	}
//
//	func NewController(repo *models.SystemConfigRepository, db *gorm.DB) *Controller {
//		return &Controller{
//			systemConfigRepo: repo,
//			baseDB:           db,
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
