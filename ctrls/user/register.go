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
		response.FailClient(c, "参数校验失败: "+err.Error())
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

	// 在插入前，先查一下当前系统里有没有用户
	userCount, err := ac.userRepo.CountUser()
	if err != nil {
		response.FailServer(c, "系统繁忙，读取配置失败")
		return
	}

	// 如果用户数为 0，说明这个人就是未来的 ID 1，给 admin 权限
	groupName := ""
	if userCount == 0 {
		groupName = "admin"
	} else {
		groupName = "user"
	}

	nowTime := time.Now()
	newUser := models.User{
		Password:  string(hashedPassword), // 存入密文
		Name:      input.Name,
		Email:     input.Email,
		Group:     groupName,
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
