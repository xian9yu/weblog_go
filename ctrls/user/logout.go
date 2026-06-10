package user

import (
	"weblog/utils"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// Logout 退出
func (ac *Controller) Logout(c *gin.Context) {
	// 调用公共函数，把当前 Token 一键送进小黑屋
	utils.Blacklist.BlockCurrentToken(c)

	response.OkMsg(c, "登出成功")
}
