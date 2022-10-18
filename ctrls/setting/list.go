package setting

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils/response"
)

func List(c *gin.Context) {

	var setting models.Setting
	total, list, err := setting.List()
	if err != nil {
		response.Error(c, "获取列表失败", err)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"total": total,
		"list":  list,
	})
}
