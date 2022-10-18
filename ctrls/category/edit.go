package category

import (
	"github.com/gin-gonic/gin"
	"time"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

func Edit(c *gin.Context) {
	id := utils.AnyToInt(c.PostForm("id"))
	name := c.PostForm("name")
	cType := c.PostForm("type")
	remark := c.PostForm("remark")
	// 判断分类名称是否存在
	count := models.CountCategoryByAny("name", name)
	if count > 0 {
		response.Error(c, "分类已存在", nil)
		c.Abort()
		return
	}

	category := models.Category{
		Id:          uint64(id),
		Name:        name,
		Type:        cType,
		Remark:      remark,
		UpdatedTime: uint64(time.Now().Unix()),
	}

	rowsAffected, err := category.Edit()

	if rowsAffected > 0 && err == nil {
		response.Success(c, gin.H{
			"message":       "保存成功",
			"rows_affected": rowsAffected,
		})
		c.Abort()
		return
	}

	response.Error(c, "保存失败, 请稍后再试", err)
}
