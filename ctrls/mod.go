package ctrls

import (
	"weblog/ctrls/article"
	"weblog/ctrls/systemConfig"
	"weblog/ctrls/user"
	"weblog/models"
)

//	type ArticleController struct {
//		// 注入 Repository 结构体，不依赖全局变量，也不直接操作 GORM
//		Repo *models.ArticleRepository
//	}
//
//	func NewArticleController(repo *models.ArticleRepository) *ArticleController {
//		return &ArticleController{Repo: repo}
//	}
type Controllers struct {
	Article      *article.Controller
	User         *user.Controller
	SystemConfig *systemConfig.Controller
	//Comment *CommentController
}

// NewControllers 注入刚才的 Repositories 大管家
func NewControllers(repos *models.Repositories) *Controllers {
	return &Controllers{
		Article:      article.NewController(repos.Article, repos.User, repos.Database),
		User:         user.NewController(repos.Article, repos.User, repos.Database),
		SystemConfig: systemConfig.NewController(repos.SystemConfig, repos.Database),
	}
}

//type UserController struct {
//	repo *models.UserRepository // 只拿自己需要的 repo
//}
//
//func NewUserController(repo *models.UserRepository) *UserController {
//	return &UserController{repo: repo}
//}
//
//type ArticleController struct {
//	repo *models.ArticleRepository // 只拿自己需要的 repo
//}
//
//func NewArticleController(repo *models.ArticleRepository) *ArticleController {
//	return &ArticleController{repo: repo}
//}
