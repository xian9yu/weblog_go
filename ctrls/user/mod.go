package user

import (
	"weblog/models"

	"gorm.io/gorm"
)

type Controller struct {
	//articleRepo *models.ArticleRepository
	baseDB   *gorm.DB
	userRepo *models.UserRepository
}

func NewController(repo *models.ArticleRepository, ur *models.UserRepository, db *gorm.DB) *Controller {
	return &Controller{
		//articleRepo: repo,
		baseDB:   db,
		userRepo: ur,
	}
}
