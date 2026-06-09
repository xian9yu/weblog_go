package systemConfig

import (
	"errors"
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetConfigById 通过ID获取完整配置详情（管理员专属）
// 路由注册：r.GET("/api/v1/admin/config/detail/:id", ac.GetConfigById)
func (ac *Controller) GetConfigById(c *gin.Context) {
	var input dto.ConfigGetByIdInput
	if err := c.ShouldBindUri(&input); err != nil {
		response.FailClient(c, "无效的配置ID")
		return
	}

	config, err := ac.systemConfigRepo.GetConfigById(input.ID)
	if err != nil {
		// 如果数据库里根本没有这个 ID
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailClient(c, "该配置项不存在或已被删除")
			return
		}
		// 其他数据库崩坏或网络异常
		response.FailServer(c, "服务器内部获取配置详情失败")
		return
	}

	// 组装出参 DTO，过滤或规范化返回的 JSON 结构
	output := dto.ConfigDetailOutput{
		Id:        config.Id,
		Key:       config.Key,
		Value:     config.Value,
		Remark:    config.Remark,
		CreatedAt: config.CreatedAt,
		UpdatedAt: config.UpdatedAt,
	}

	// 4. 丝滑响应
	response.Ok(c, output)
}

// GetConfigByKey 快速获取单个配置的 value
// 路由：r.GET("/api/v1/admin/config/:key", ac.GetConfigByKey)
func (ac *Controller) GetConfigByKey(c *gin.Context) {
	var input dto.ConfigGetByKeyInput

	if err := c.ShouldBindUri(&input); err != nil {
		response.FailClient(c, "请求的配置项格式不正确")
		return
	}

	value, err := ac.systemConfigRepo.GetValueByKey(input.Key)
	if err != nil {
		response.FailServer(c, "获取配置失败")
		return
	}

	response.Ok(c, gin.H{"value": value})
}
