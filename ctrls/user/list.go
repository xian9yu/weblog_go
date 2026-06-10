package user

import (
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// List 查询 用户列表
func (ac *Controller) List(c *gin.Context) {
	tokenGroup, _ := c.Get("group")
	if tokenGroup != "admin" {
		// 非管理员访问，返回 403 Forbidden
		response.FailForbidden(c, "权限不足，仅管理员可查看用户列表")
		c.Abort()
		return
	}

	//使用 DTO 一键绑定并过滤 Query 参数
	var input dto.UserListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "参数校验失败: "+err.Error())
		c.Abort()
		return
	}

	// 如果前端没传 page_no，或者传了 0，默认第 1 页
	if input.PageNo <= 0 {
		input.PageNo = 1
	}

	// 如果前端没传 page_size，或者传了 0，默认每页 10 条
	if input.PageSize <= 0 {
		input.PageSize = 10
	}

	// 3. ✨ 优雅调用：交由 Repo 仓储层去处理繁琐的数据库分页
	total, list, err := ac.repos.User.GetList(input)
	if err != nil {
		response.FailServer(c, "获取用户列表失败")
		c.Abort()
		return
	}

	// 4. 统一成功规范返回
	response.Ok(c, gin.H{
		"total":     total,
		"page_no":   input.PageNo,
		"page_size": input.PageSize,
		"list":      list,
	})
}
