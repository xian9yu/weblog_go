package article

import (
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// GetDetailsById article 详情
func (ac *Controller) GetDetailsById(c *gin.Context) {
	var input dto.ArticleDetailInput
	if err := c.ShouldBindQuery(&input); err != nil {
		response.FailClient(c, "文章查询参数格式不正确")
		return
	}
	article, err := ac.articleRepo.GetArticleDetailsById(input.ID)
	if err != nil || article == nil {
		response.FailClient(c, "文章不存在或已被删除")
		c.Abort()
		return
	}

	//article.Id
	user, err := ac.userRepo.GetDetailsById(article.Id)
	if err != nil || user == nil {
		response.FailClient(c, "用户不存在或已被删除")
		c.Abort()
		return
	}
	output := dto.ArticleDetailResponse{
		ID:          article.Id,
		Title:       article.Title,
		Content:     article.Content,
		AuthorName:  user.Name,
		Status:      article.Status,
		CreatedTime: article.CreatedTime,
		UpdatedTime: article.UpdatedTime,
	}

	response.Ok(c, gin.H{
		"details": output,
	})
}
