package systemConfig

import (
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// GetConfigList 获取系统配置项分页列表（管理员专属）
// 适用于路由：GET /api/v1/admin/config
func (ac *Controller) GetConfigList(c *gin.Context) {
	// 绑定并强校验 URL 查询参数（?page=1&page_size=10&keyword=site）
	var input dto.ConfigListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		response.FailClient(c, err.Error())
		return
	}

	// 调用 Repo 层执行高性能、防污染的分页查询
	modelsList, total, err := ac.systemConfigRepo.GetConfigList(input)
	if err != nil {
		response.FailServer(c, "获取配置列表失败")
		return
	}

	// 组装出参 DTO 数组，过滤不必要的底层细节，规范 JSON 字段
	items := make([]dto.ConfigInfoOutput, 0, len(modelsList))
	for _, config := range modelsList {
		items = append(items, dto.ConfigInfoOutput{
			Id:        config.Id,
			Key:       config.Key,
			Value:     config.Value,
			Remark:    config.Remark,
			CreatedAt: config.CreatedAt,
			UpdatedAt: config.UpdatedAt,
		})
	}

	if len(items) == 0 {
		items = []dto.ConfigInfoOutput{}
	}

	output := dto.ConfigListOutput{
		List:  items,
		Total: total,
	}

	response.Ok(c, output)
}
