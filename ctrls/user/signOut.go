package user

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

// 退出
func SignOut(c *gin.Context) {
	id := c.PostForm("user_id")
	userId := uint64(utils.AnyToInt(id))

	// 获取 token 详情
	tokenUserId, _ := c.Get("user_id")
	group, _ := c.Get("group")
	// 只有 (本人/管理员) 才能操作
	if userId == uint64(utils.AnyToInt(tokenUserId)) || group == "admin" {

		auth := models.Auth{
			UserId: userId,
		}

		// 登出
		rowsAffected := auth.Delete()
		if rowsAffected == 1 {
			response.Success(c, gin.H{
				"message": "登出成功",
			})
			c.Abort()
			return
		}

		response.Error(c, "登出失败", nil)
		c.Abort()
		return
	}

	response.PageNotFound(c)
}
