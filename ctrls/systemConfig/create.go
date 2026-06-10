package systemConfig

import (
	"errors"
	"regexp"
	"weblog/dto"
	"weblog/models"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 预编译一个正则：只允许英文、数字、下划线和中划线
var keyRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// CreateConfig 新增配置项（管理员专属）
// 路由注册：r.POST("/api/v1/admin/config", ac.CreateConfig)
func (ac *Controller) CreateConfig(c *gin.Context) {
	var input dto.ConfigCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, err.Error())
		return
	}

	// 校验 Key 的命名规范
	if !keyRegex.MatchString(input.Key) {
		response.FailClient(c, "配置键名格式不正确！只允许使用英文、数字、下划线(_)或中划线(-)")
		return
	}

	// 防止 Key 重复创建（唯一性检查）
	_, err := ac.systemConfigRepo.GetValueByKey(input.Key)
	if err == nil {
		response.FailClient(c, "该配置 key 已存在，请勿重复创建")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		response.FailServer(c, "检测配置键名冲突时发生异常")
		return
	}

	newConfig := models.SystemConfig{
		Key:    input.Key,
		Value:  input.Value,
		Remark: input.Remark,
	}

	if err := ac.systemConfigRepo.CreateConfig(&newConfig); err != nil {
		response.FailServer(c, "新建配置项失败，数据库写入异常")
		return
	}

	response.OkMsg(c, "新配置项写入成功！")
}
