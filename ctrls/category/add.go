package category

import (
	"github.com/gin-gonic/gin"
	"time"
	"weblog/models"
	"weblog/utils/response"
)

func Add(c *gin.Context) {
	name := c.PostForm("name")
	cType := c.PostForm("type")
	remark := c.PostForm("remark")

	if name == "" {
		response.Error(c, "name 不能为空", nil)
		c.Abort()
		return
	}
	if cType == "" {
		response.Error(c, "type 不能为空", nil)
		c.Abort()
		return
	}

	// 判断分类名称是否存在
	count := models.CountCategoryByAny("name", name)
	if count > 0 {
		response.Error(c, "分类已存在", nil)
		c.Abort()
		return
	}

	nowTime := uint64(time.Now().Unix())
	category := models.Category{
		Name:        name,
		Type:        cType,
		Remark:      remark,
		CreatedTime: nowTime,
		UpdatedTime: nowTime,
	}

	// 添加分类
	categoryId, rowsAffected, err := category.Add()
	if rowsAffected > 0 && err == nil {
		response.Success(c, gin.H{
			"message":       "保存成功",
			"rows_affected": rowsAffected,
			"categoryId":    categoryId,
		})
		c.Abort()
		return
	}

	response.Error(c, "保存失败, 请稍后再试", err)
}
