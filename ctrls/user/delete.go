package user

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

func Delete(c *gin.Context) {
	id := c.PostForm("id")
	uId := uint64(utils.AnyToInt(id))

	// 获取 token 详情
	userId, _ := c.Get("user_id")
	group, _ := c.Get("group")
	// 只有 (本人/管理员) 才能操作
	if uId == uint64(utils.AnyToInt(userId)) || group == "admin" {
		var user models.User
		rowsAffected := user.Delete(uId)
		if rowsAffected == 1 {
			response.Success(c, gin.H{
				"message": "删除成功",
			})
			c.Abort()
			return
		}

		response.Error(c, "删除失败, 请稍后再试", nil)
		c.Abort()
		return
	}

	response.PageNotFound(c)
}
