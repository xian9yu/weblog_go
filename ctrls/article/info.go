package article

import (
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// GetInfoById article 详情
func (ac *Controller) GetInfoById(c *gin.Context) {
	var input dto.ArticleInfoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "文章查询参数格式不正确")
		c.Abort()
		return
	}
	article, err := ac.repos.Article.GetArticleInfosById(input.ID)
	if err != nil || article == nil {
		response.FailClient(c, "文章不存在或已被删除")
		c.Abort()
		return
	}

	user, err := ac.repos.User.GetGetUserInfoById(article.UserId)

	if err != nil || user == nil {
		response.FailClient(c, "文章不存在或已被删除")
		c.Abort()
		return
	}
	output := dto.ArticleInfoResponse{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		CategoryName: article.CategoryName,
		AuthorName:   user.Name,
		Status:       article.Status,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}

	response.Ok(c, gin.H{
		"details": output,
	})
}
