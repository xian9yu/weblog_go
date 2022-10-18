package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

func Delete(c *gin.Context) {
	idStr := c.PostForm("id")
	id := uint64(utils.AnyToInt(idStr))

	// 判断分类名称是否存在
	count := models.CountCategoryByAny("id", idStr)
	if count < 1 {
		response.Error(c, "分类 id 不存在", nil)
		c.Abort()
		return
	}

	// 事务开始
	var db = models.MySQL()
	err := db.Transaction(func(tx *gorm.DB) error {
		// 删除分类
		var category models.Category
		err := category.TxDelete(tx, id)
		if err != nil {
			return err
		}

		// 通过分类 id 删除 关联文章的 id
		var links models.Links
		links.CategoryId = id
		err = links.TxDeleteByCategoryId(tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		response.Error(c, "删除失败, 请重试", nil)
	}

	response.Error(c, "删除成功", nil)

}
