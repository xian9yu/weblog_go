package systemConfig

import (
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// DeleteConfig 批量/单个删除系统配置项（管理员专属）
// 适用于路由：DELETE /api/v1/admin/config
func (ac *Controller) DeleteConfig(c *gin.Context) {
	var input dto.ConfigDeleteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "参数错误: ID列表不能为空")
		c.Abort()
		return
	}

	err := ac.repos.SystemConfig.DeleteConfigs(input.Ids)
	if err != nil {
		response.FailServer(c, "删除配置失败")
		c.Abort()
		return
	}

	// 3. 丝滑返回成功
	response.Ok(c, "删除成功")
}
