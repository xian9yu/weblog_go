package links

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils/response"
)

func ListWithArticle(c *gin.Context) {

	articleId := c.Query("article_id") // 分类id

	if articleId == "" {
		response.Error(c, "article id is null", nil)
		c.Abort()
		return
	}
	var articleIds []int
	// 把 string 转为 []string
	_ = json.Unmarshal([]byte(articleId), &articleIds)

	var links models.Links
	list, err := links.GetList(articleIds)
	if err != nil {
		response.Error(c, "获取列表失败", err)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"article_id": "",
		"list":       list,
	})
}
