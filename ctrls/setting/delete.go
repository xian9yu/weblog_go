package setting

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

func Delete(c *gin.Context) {
	id := c.PostForm("id")

	var setting models.Setting
	rowsAffected := setting.Delete(uint64(utils.AnyToInt(id)))

	if rowsAffected == 1 {
		response.Success(c, gin.H{
			"message": "删除成功",
		})
		c.Abort()
		return
	}

	response.Error(c, "删除失败, 请稍后再试", nil)

}
