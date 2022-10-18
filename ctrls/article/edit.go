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

func Edit(c *gin.Context) {
	id := c.PostForm("id")
	title := c.PostForm("title")
	content := c.PostForm("content")
	allowView := c.PostForm("allow_view")
	categoryId := c.PostForm("category_id")

	count := models.CountArticleByAny("title", title)
	if count > 0 {
		response.Error(c, "title 已存在", nil)
		c.Abort()
		return
	}

	// 获取 token 详情
	userId, _ := c.Get("user_id")
	group, _ := c.Get("group")
	// 只有 (本人/管理员) 才能操作
	if uint64(utils.AnyToInt(id)) == uint64(utils.AnyToInt(userId)) || group == "admin" {
		// 事务开始
		var db = models.MySQL()
		err := db.Transaction(func(tx *gorm.DB) error {
			// 删除所有 article 对应的 category id
			var links models.Links
			links.ArticleId = uint64(utils.AnyToInt(id))
			err := links.TxDeleteByArticleId(tx)
			if err == nil {
				// 修改 article 内容
				article := models.Article{
					Id:          links.ArticleId,
					Title:       title,
					Content:     content,
					AllowView:   allowView,
					UpdatedTime: uint64(time.Now().Unix()),
				}
				affected, errE := article.TxEdit(tx)
				if affected < 1 || errE != nil {
					return errE
				}

				var (
					categoryIds []int64
					sliceLinks  []models.Links
				)
				// 把 string 转为 []string
				_ = json.Unmarshal([]byte(categoryId), &categoryIds)
				if categoryIds == nil {
					categoryIds = append(categoryIds, 1) // 默认分类
				}

				for _, cId := range categoryIds {
					for j := 0; j < len(categoryIds); j++ {
						links = models.Links{
							ArticleId:   links.ArticleId,
							CategoryId:  uint64(cId),
							UpdatedTime: uint64(time.Now().Unix()),
						}
					}
					sliceLinks = append(sliceLinks, links)
				}

				// 添加文章分类 links
				rowsAffected, err := links.TxAdd(tx, sliceLinks)
				if rowsAffected < 1 || err != nil {
					return err
				}

				return nil
			}

			return err
		})
		if err != nil {
			response.Error(c, "article 保存失败, 请稍后再试", err)
			c.Abort()
			return
		}

		response.Success(c, gin.H{
			"message": "保存成功",
		})
		c.Abort()
		return
	}

	response.PageNotFound(c)
}
