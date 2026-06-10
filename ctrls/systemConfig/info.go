package systemConfig

import (
	"errors"
	"fmt"
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetBatchValues 批量获取配置项的值
// POST /api/v1/config/batch
func (ac *Controller) GetBatchValues(c *gin.Context) {
	var input dto.ConfigBatchGetByKeysInput

	fmt.Println(input.Keys)
	//  keys 是个数组，用 POST + JSON 传参最稳妥
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "参数错误: keys 列表不能为空")
		c.Abort()
		return
	}

	// 批量全家桶查询
	configMap, err := ac.repos.SystemConfig.GetValuesByKeys(input.Keys)
	if err != nil {
		response.FailServer(c, "批量获取配置失败")
		c.Abort()
		return
	}

	response.Ok(c, configMap)
}

// GetConfigById 通过ID获取完整配置详情
// 路由注册：r.GET("/api/v1/admin/config/info", ac.GetConfigById)
func (ac *Controller) GetConfigById(c *gin.Context) {
	var input dto.ConfigGetByIdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "无效的配置ID")
		c.Abort()
		return
	}

	config, err := ac.repos.SystemConfig.GetConfigById(input.ID)
	if err != nil {
		// 如果数据库里根本没有这个 ID
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailClient(c, "该配置项不存在或已被删除")
			c.Abort()
			return
		}
		// 其他数据库崩坏或网络异常
		response.FailServer(c, "服务器内部获取配置详情失败")
		c.Abort()
		return
	}

	// 组装出参 DTO，过滤或规范化返回的 JSON 结构
	output := dto.ConfigInfoOutput{
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
		c.Abort()
		return
	}

	value, err := ac.repos.SystemConfig.GetValueByKey(input.Key)
	if err != nil {
		response.FailServer(c, "获取配置失败")
		c.Abort()
		return
	}

	response.Ok(c, gin.H{"value": value})
}
