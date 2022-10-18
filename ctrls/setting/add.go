package setting

import (
	"github.com/gin-gonic/gin"
	"time"
	"weblog/models"
	"weblog/utils/response"
)

func Add(c *gin.Context) {
	name := c.PostForm("name")
	value := c.PostForm("value")
	remark := c.PostForm("remark")
	// 参数是否为空
	if name == "" {
		response.Error(c, "name 不能为空", nil)
		c.Abort()
		return
	}
	if value == "" {
		response.Error(c, "value 不能为空", nil)
		c.Abort()
		return
	}

	// 判断 setting name 是否存在
	nameCount := models.CountSettingByAny("name", name)
	if nameCount > 0 {
		response.Error(c, "设置已存在", nil)
		c.Abort()
		return
	}

	nowTime := uint64(time.Now().Unix())
	setting := &models.Setting{
		Name:        name,
		Value:       value,
		Remark:      remark,
		CreatedTime: nowTime,
		UpdatedTime: nowTime,
	}

	id, rowsAffected, err := setting.Add()
	if rowsAffected == 1 && id > 0 {
		response.Success(c, gin.H{
			"message":    "保存成功",
			"setting_id": id,
		})
		c.Abort()
		return
	}

	response.Error(c, "保存失败, 请稍后再试", err)
}
