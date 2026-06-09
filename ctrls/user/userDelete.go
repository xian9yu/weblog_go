package user

import (
	"weblog/dto"
	"weblog/utils"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// Unregister 自主注销账号
// 路由注册：POST /api/v1/admin/user/unregister （必须挂载在 JWT 中建件后面）
func (ac *Controller) Unregister(c *gin.Context) {
	currentUserId, _ := c.Get("userId")
	uid, _ := currentUserId.(uint64)

	var input dto.UserUnregisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, err.Error())
		return
	}

	// 调用 Repo 层进行注销
	success, err := ac.userRepo.UnregisterSelf(uid, input.Password)
	if err != nil {
		response.FailServer(c, "注销失败")
		return
	}
	if !success {
		response.FailClient(c, "密码错误")
		return
	}

	//  顺手把被注销者的 Token 拉黑
	utils.Blacklist.BlockCurrentToken(c)
	response.OkMsg(c, "您的账号已成功注销，江湖路远，有缘再见！")
}
