package links

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

func Add(c *gin.Context) {
	articleId := c.PostForm("article_id")
	category := c.PostForm("category_id")

	var (
		categoryIds []int64
		sliceLinks  []models.Links
	)

	// 把 string 转为 []string
	_ = json.Unmarshal([]byte(category), &categoryIds)
	if categoryIds == nil {
		categoryIds = append(categoryIds, 1)
	}

	if articleId == "" {
		response.Error(c, "article_id 不能为空", nil)
		c.Abort()
		return
	}

	if category == "" {
		response.Error(c, "category_id 不能为空", nil)
		c.Abort()
		return
	}

	links := models.Links{
		ArticleId:   uint64(utils.AnyToInt(articleId)),
		CategoryId:  uint64(utils.AnyToInt(category)),
		UpdatedTime: uint64(time.Now().Unix()),
	}

	count := links.CountLinksById()
	if count > 0 {
		response.Success(c, gin.H{})

		c.Abort()
		return
	}

	var l models.Links
	for _, cId := range categoryIds {
		for j := 0; j < len(categoryIds); j++ {
			l = models.Links{
				ArticleId:   uint64(utils.AnyToInt(articleId)),
				CategoryId:  uint64(cId),
				UpdatedTime: uint64(time.Now().Unix()),
			}
		}
		sliceLinks = append(sliceLinks, links)
	}
	// 添加分类
	rowsAffected, err := l.Add(sliceLinks)
	if rowsAffected > 0 && err == nil {
		response.Success(c, gin.H{
			"message":       "链接添加成功",
			"rows_affected": rowsAffected,
		})
		c.Abort()
		return
	}

	response.Error(c, "链接添加失败", err)
}
