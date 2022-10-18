package article

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

// GetDetailsById article 详情
func GetDetailsById(c *gin.Context) {
	articleId := c.Param("id")

	result, err := models.GetArticleDetailsByAny("id", utils.AnyToInt(articleId))
	if err != nil {
		response.Error(c, "获取内容失败", err)
		c.Abort()
		return
	}
	if result.Id < 1 {
		response.Error(c, "文章不存在", nil)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"details": result,
	})
}

// GetDetailsByTitle article 详情
func GetDetailsByTitle(c *gin.Context) {
	articleId := c.Param("title")

	result, err := models.GetArticleDetailsByAny("title", articleId)
	if err != nil {
		response.Error(c, "获取内容失败", err)
		c.Abort()
		return
	}

	if result.Id < 1 {
		response.Error(c, "文章不存在", nil)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"details": result,
	})
}
