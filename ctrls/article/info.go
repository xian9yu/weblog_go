package article

import (
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// GetInfoById article 详情
func (ac *Controller) GetInfoById(c *gin.Context) {
	var input dto.ArticleInfoInput
	if err := c.ShouldBindQuery(&input); err != nil {
		response.FailClient(c, "文章查询参数格式不正确")
		return
	}
	article, err := ac.articleRepo.GetArticleInfosById(input.ID)
	if err != nil || article == nil {
		response.FailClient(c, "文章不存在或已被删除")
		c.Abort()
		return
	}

	//article.Id
	user, err := ac.userRepo.GetGetUserInfoById(article.ID)
	if err != nil || user == nil {
		response.FailClient(c, "用户不存在或已被删除")
		c.Abort()
		return
	}
	output := dto.ArticleInfoResponse{
		ID:         article.ID,
		Title:      article.Title,
		Content:    article.Content,
		AuthorName: user.Name,
		Status:     article.Status,
		CreatedAt:  article.CreatedAt,
		UpdatedAt:  article.UpdatedAt,
	}

	response.Ok(c, gin.H{
		"details": output,
	})
}
