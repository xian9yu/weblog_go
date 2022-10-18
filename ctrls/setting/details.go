package setting

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

// 详情
func GetDetailsById(c *gin.Context) {
	id := c.Param("id")

	result, err := models.GetSettingDetailsByAny("id", utils.AnyToInt(id))
	if err != nil {
		response.Error(c, "获取 setting 失败", err)
		c.Abort()
		return
	}
	if result.Id < 1 {
		response.Error(c, "数据不存在", nil)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"details": result,
	})
}

// 详情
func GetDetailsByName(c *gin.Context) {
	name := c.Param("name")

	result, err := models.GetSettingDetailsByAny("name", name)
	if err != nil {
		response.Error(c, "获取 setting 失败", err)
		c.Abort()
		return
	}

	if result.Id < 1 {
		response.Error(c, "数据不存在", nil)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"details": result,
	})
}
