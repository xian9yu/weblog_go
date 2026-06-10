package user

import (
	"errors"
	"strconv"
	"time"
	"weblog/dto"
	"weblog/utils"

	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login 登录
func (ac *Controller) Login(c *gin.Context) {
	var input dto.UserLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "参数校验失败")
		return
	}

	// 获取用户信息
	user, err := ac.repos.User.LoginByEmail(input.Email)
	if err != nil {
		// 为了防止黑客暴力破解，无论“用户不存在”还是“密码错误”，统一返回模糊的提示
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailClient(c, "账号或密码错误")
			return
		}
		response.FailServer(c, "系统繁忙，登录失败: "+err.Error())
		return
	}

	// 检查账号状态 1启用 2禁用
	if user.Status != 1 {
		response.FailForbidden(c, "该账号停止使用")
		return
	}

	// 比对密码哈希值
	// input.Password 是前端传来的明文密码（如 "123456"）
	// user.Password 是数据库里存的密文哈希（如 "$2a$10$X..."）
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// 密码对不上，报错拦截
		response.FailClient(c, "账号或密码错误")
		return
	}

	// 生成 token
	token, _ := utils.GenerateToken(strconv.FormatUint(user.ID, 10), user.Name, user.Email, user.Group, 8*time.Hour)

	response.Ok(c, dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
