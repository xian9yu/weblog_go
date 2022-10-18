package links

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

func Delete(c *gin.Context) {
	articleId := c.PostForm("article_id")
	categoryId := c.PostForm("category_id")

	var links models.Links
	rowsAffected := links.Delete(uint64(utils.AnyToInt(articleId)), uint64(utils.AnyToInt(categoryId)))

	if rowsAffected == 1 {
		response.Success(c, gin.H{
			"message": "删除成功",
		})

		c.Abort()
		return
	}

	response.Error(c, "article not found", nil)

}
