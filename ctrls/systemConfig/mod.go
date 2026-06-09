package systemConfig

import (
	"weblog/models"

	"gorm.io/gorm"
)

type Controller struct {
	systemConfigRepo *models.SystemConfigRepository // 只拿自己需要的 repo
	baseDB           *gorm.DB
}

func NewController(repo *models.SystemConfigRepository, db *gorm.DB) *Controller {
	return &Controller{
		systemConfigRepo: repo,
		baseDB:           db,
	}
}
