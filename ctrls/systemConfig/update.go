package systemConfig

import (
	"errors"
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateConfigByKey 修改配置项（管理员专属）
// 路由注册：r.PUT("/api/v1/admin/config", ac.UpdateConfigByKey)
func (ac *Controller) UpdateConfigByKey(c *gin.Context) {
	var input dto.ConfigUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, err.Error())
		return
	}

	// 先检查这个 Key 到底在不在数据库里
	_, err := ac.systemConfigRepo.GetValueByKey(input.Key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailClient(c, "该配置项不存在，无法进行修改")
			return
		}
		response.FailServer(c, "检查配置项完整性失败")
		return
	}

	if err := ac.systemConfigRepo.UpdateValueByKey(input.Key, input.Value, input.Remark); err != nil {
		response.FailServer(c, "保存配置失败，数据库写入异常")
		return
	}

	// 4. 丝滑返回
	response.OkMsg(c, "配置项已成功更新")
}
