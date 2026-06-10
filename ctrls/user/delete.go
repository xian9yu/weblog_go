package user

import (
	"fmt"
	"weblog/dto"
	"weblog/utils"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// DeleteAccount 自主注销账号
// 路由注册：POST /api/v1/admin/user/delete （必须挂载在 JWT 中建件后面）
func (ac *Controller) DeleteAccount(c *gin.Context) {
	currentUserId := c.GetUint64("user_id")
	// 如果是唯一的超级管理员，拒绝自残行为！
	if currentUserId == 1 {
		response.FailClient(c, "为了系统安全，主管理员账号禁止被删除！")
		return
	}

	var input dto.UserUnregisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, err.Error())
		return
	}
	fmt.Println(input, currentUserId)

	// 调用 Repo 层进行注销
	success, err := ac.repos.User.DeleteAccount(currentUserId, input.Password)
	if err != nil {
		response.FailServer(c, "注销失败，系统繁忙")
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
