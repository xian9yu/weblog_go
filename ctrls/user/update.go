package user

import (
	"time"
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// UpdateProfile 作者修改个人资料（不含密码）
// 路由：PUT /api/v1/admin/user/profile
func (ac *Controller) UpdateProfile(c *gin.Context) {
	currentUserId, _ := c.Get("userId")

	// 绑定 DTO
	var input dto.UserUpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, err.Error())
		c.Abort()
		return
	}

	// DTO 转 Map，并进行严格的安全过滤
	userData := map[string]any{
		"name":       input.Name,
		"email":      input.Email,
		"updated_at": time.Now(),
	}

	// 调用 Repo 执行更新
	rowsAffected, err := ac.repos.User.UpdateProfile(currentUserId.(uint64), userData)
	if err != nil {
		response.FailServer(c, "修改资料失败")
		c.Abort()
		return
	}

	if rowsAffected == 0 {
		response.FailClient(c, "资料未做任何变动")
		c.Abort()
		return
	}

	response.OkMsg(c, "个人资料修改成功")
}

// UpdatePassword 修改密码接口
// 路由：PUT /api/v1/admin/user/password
func (ac *Controller) UpdatePassword(c *gin.Context) {
	currentUserId, _ := c.Get("userId")
	uid, _ := currentUserId.(uint64)

	var input dto.UserUpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, err.Error())
		c.Abort()
		return
	}

	// 调用 Repo 层高内聚的方法
	ok, err := ac.repos.User.UpdatePassword(uid, input.OldPassword, input.NewPassword)
	if err != nil {
		response.FailServer(c, "修改密码系统故障")
		c.Abort()
		return
	}

	if !ok {
		response.FailClient(c, "原密码输入错误，请重新输入")
		c.Abort()
		return
	}

	response.OkMsg(c, "密码修改成功，请牢记新密码")
}

// AdminResetUserPassword 管理员后台强制修改他人密码
// 路由：PUT /api/v1/admin/user/reset-password （管理员专属路由）
func (ac *Controller) AdminResetUserPassword(c *gin.Context) {
	// 严格的管理员权限校验
	tokenGroup, _ := c.Get("group")
	if tokenGroup != "admin" {
		response.FailForbidden(c, "权限不足，仅管理员可重置他人密码")
		c.Abort()
		return
	}

	// 绑定 DTO 参数
	var input dto.AdminResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, err.Error())
		c.Abort()
		return
	}

	// 防止管理员“误操作”或者“降权操作”
	// 如果管理员不小心在后台输入了自己的 UserID，虽然能成功，但通常建议提示他走正常的“修改密码”流程
	currentUserId, _ := c.Get("userId")
	uid, _ := currentUserId.(uint64)
	if input.TargetUserId == uid {
		response.FailClient(c, "修改自身密码请走『修改个人密码』接口，需验旧密码")
		c.Abort()
		return
	}

	// 调用 Repo 层的强制更新功能
	err := ac.repos.User.ForceUpdatePassword(input.TargetUserId, input.NewPassword)
	if err != nil {
		response.FailServer(c, "重置密码失败，数据库写入异常")
		c.Abort()
		return
	}

	response.OkMsg(c, "该用户密码已成功重置")
}
