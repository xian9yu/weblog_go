package user

import (
	"weblog/models"

	"gorm.io/gorm"
)

//	type Controller struct {
//		//articleRepo *models.ArticleRepository
//		baseDB           *gorm.DB
//		userRepo         *models.UserRepository
//		systemConfigRepo *models.SystemConfigRepository
//	}
//
//	func NewController(scr *models.SystemConfigRepository, ur *models.UserRepository, db *gorm.DB) *Controller {
//		return &Controller{
//			//articleRepo: repo,
//			baseDB:           db,
//			userRepo:         ur,
//			systemConfigRepo: scr,
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
