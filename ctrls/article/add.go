package article

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

func Add(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	allowView := c.PostForm("allow_view")
	category := c.PostForm("category_id")

	if title == "" {
		response.Error(c, "title 不能为空", nil)
		c.Abort()
		return
	}
	if content == "" {
		response.Error(c, "content 不能为空", nil)
		c.Abort()
		return
	}

	count := models.CountArticleByAny("title", title)
	if count > 0 {
		response.Error(c, "title 已存在", nil)
		c.Abort()
		return
	}

	// 获取 token 详情
	userId, _ := c.Get("user_id")

	nowTime := uint64(time.Now().Unix())
	article := &models.Article{
		Title:       title,
		Content:     content,
		UserId:      uint64(utils.AnyToInt(userId)),
		AllowView:   allowView, // y 显示 | n 不显示
		CreatedTime: nowTime,
		UpdatedTime: nowTime,
	}

	// 事务开始
	var db = models.MySQL()
	err := db.Transaction(func(tx *gorm.DB) error {
		// 新增 article
		articleId, rowsAffected, err := article.TxAdd(tx)
		if rowsAffected < 1 || err != nil {
			return err
		}

		var (
			categoryIds []int64
			sliceLinks  []models.Links
		)

		// 把 string 转为 []string
		_ = json.Unmarshal([]byte(category), &categoryIds)
		if categoryIds == nil {
			categoryIds = append(categoryIds, 1)
		}

		for _, cId := range categoryIds {
			var links models.Links
			for j := 0; j < len(categoryIds); j++ {
				links = models.Links{
					ArticleId:   articleId,
					CategoryId:  uint64(cId),
					UpdatedTime: nowTime,
				}
			}
			sliceLinks = append(sliceLinks, links)
		}

		// 添加分类 links
		var links models.Links
		rowsAffected, err = links.TxAdd(tx, sliceLinks)
		if rowsAffected < 1 || err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		response.Error(c, "保存失败, 请稍后再试", err)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"message": "保存成功",
	})
}
