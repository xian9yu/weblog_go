package user

import (
	"time"
	"weblog/dto"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register 注册
func (ac *Controller) Register(c *gin.Context) {
	var input dto.UserRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "参数校验失败")
		return
	}

	// 判断是否允许注册
	//systemConfig, er := models.GetSettingDetailsByAny("name", "allow_signup")
	//if er != nil {
	//	response.Error(c, "注册失败, 请重试", er)
	//	c.Abort()
	//	return
	//}
	//if systemConfig.Value != "y" {
	//	response.Error(c, "未开启注册", nil)
	//	c.Abort()
	//	return
	//}

	// 判断 邮箱 存在
	isEmailExist, err := ac.userRepo.CheckEmailExist(input.Email)
	if err != nil {
		response.FailClient(c, "参数校验失败，请检查输入格式")
		return
	}
	if isEmailExist {
		response.FailClient(c, "该邮箱已被其他账号绑定")
		return
	}

	// 将明文密码哈希化，防范数据库泄露风险
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		response.FailServer(c, "密码加密失败")
		return
	}

	nowTime := time.Now()
	newUser := models.User{
		Password:  string(hashedPassword), // 存入密文
		Name:      input.Name,
		Email:     input.Email,
		Status:    1, // 默认允许登录
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	userId, rowsAffected, err := ac.userRepo.Create(&newUser)
	if err != nil || rowsAffected == 0 {
		response.FailServer(c, "系统繁忙，创建用户失败")
		return
	}

	response.OkMsg(c, "新用户注册成功，ID为: "+utils.AnyToString(userId))
}
