package user

import (
	"errors"
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetDetailsById 获取详情
func (ac *Controller) GetDetailsById(c *gin.Context) {
	var input dto.UserDetailInput

	// 针对 uri:"id" 标签，必须使用 ShouldBindUri 一键绑定和格式化
	// 路由定义如：r.GET("/api/v1/admin/user/:id", ac.GetUserDetail)
	if err := c.ShouldBindUri(&input); err != nil {
		response.FailClient(c, "用户ID格式不正确")
		return
	}

	user, err := ac.userRepo.GetUserDetailsById(input.ID)
	if err != nil {
		// 精准拦截 GORM 的查无此人错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailClient(c, "该用户不存在")
			return
		}
		response.FailServer(c, "获取用户详情失败")
		return
	}

	response.Ok(c, user)
}

// GetDetailsByEmail 获取详情
func (ac *Controller) GetDetailsByEmail(c *gin.Context) {
	var input dto.UserDetailInput
	// 绑定路径中的 :email（Gin 会自动进行 URL 解码，将 %40 还原为 @）
	if err := c.ShouldBindUri(&input); err != nil {
		response.FailClient(c, "邮箱格式不正确")
		return
	}

	user, err := ac.userRepo.GetUserDetailsByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailClient(c, "该邮箱对应的用户不存在")
			return
		}
		response.FailServer(c, "系统异常，获取详情失败")
		return
	}

	response.Ok(c, user)
}
