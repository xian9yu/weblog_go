package category

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils/response"
)

// 详情 by name
func GetDetailsByName(c *gin.Context) {
	name := c.Param("name")
	result, err := models.GetCategoryDetailsByAny("name", name)
	if err != nil {
		response.Error(c, "获取详情失败", err)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"details": result,
	})
}

// 详情 by id
func GetDetailsById(c *gin.Context) {
	id := c.Param("id")

	result, err := models.GetCategoryDetailsByAny("id", id)
	if err != nil {
		response.Error(c, "获取详情失败", err)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"details": result,
	})
}
