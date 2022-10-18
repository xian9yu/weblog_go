package setting

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
	value := c.PostForm("value")
	remark := c.PostForm("remark")

	setting := models.Setting{
		Id:          uint64(id),
		Name:        name,
		Value:       value,
		Remark:      remark,
		UpdatedTime: uint64(time.Now().Unix()),
	}

	rowsAffected, err := setting.Edit()

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
