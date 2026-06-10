package ctrls

import (
	"weblog/ctrls/article"
	"weblog/ctrls/systemConfig"
	"weblog/ctrls/user"
	"weblog/models"
)

type Controllers struct {
	Article      *article.Controller
	User         *user.Controller
	SystemConfig *systemConfig.Controller
	//Comment *CommentController
}

// NewControllers 注入 Repositories
func NewControllers(repos *models.Repositories) *Controllers {
	return &Controllers{
		Article:      article.NewController(repos),
		User:         user.NewController(repos),
		SystemConfig: systemConfig.NewController(repos),
	}
}
