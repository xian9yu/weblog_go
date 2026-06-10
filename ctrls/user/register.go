package user

import (
	"errors"
	"time"
	"weblog/dto"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register 注册
func (ac *Controller) Register(c *gin.Context) {
	var input dto.UserRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "参数校验失败: "+err.Error())
		c.Abort()
		return
	}

	// 检查是否允许注册
	canRegister, err := ac.repos.SystemConfig.GetBool("allow_registration")
	if err != nil {
		// 错误是“没找到记录”，说明是第一个用户来冷启动，放行！
		if errors.Is(err, gorm.ErrRecordNotFound) {
			canRegister = true // 找不到配置，默认允许注册
		} else {
			// 如果是其他数据库挂了的错误（比如连接断开），才真正报错拦截
			response.FailServer(c, "系统服务繁忙，请稍后再试")
			c.Abort()
			return
		}
	}
	// 如果配置明确存在，且被设置为了 false，强行拦截
	if !canRegister {
		response.FailClient(c, "当前系统已关闭注册功能")
		c.Abort()
		return
	}

	// 判断 邮箱 存在
	isEmailExist, err := ac.repos.User.CheckEmailExist(input.Email)
	if err != nil {
		response.FailClient(c, "参数校验失败，请检查输入格式")
		c.Abort()
		return
	}
	if isEmailExist {
		response.FailClient(c, "该邮箱已被其他账号绑定")
		c.Abort()
		return
	}

	// 将明文密码哈希化，防范数据库泄露风险
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		response.FailServer(c, "密码加密失败")
		c.Abort()
		return
	}

	// 准备用于拼装返回的临时变量
	var createdUserId uint64

	// 开启事务
	err = ac.baseDB.Transaction(func(tx *gorm.DB) error {
		txRepos := ac.repos.WithTx(tx)

		// 在事务内查一下当前系统里有没有用户（并发安全）
		userCount, err := txRepos.User.CountUser()
		if err != nil {
			return err // 报错则自动回滚
		}

		// 2. 如果用户数为 0，给 admin 权限
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

		// 使用处于事务状态的 txRepos.User 执行写入
		userId, rowsAffected, err := txRepos.User.Create(&newUser)
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errors.New("创建用户失败")
		}

		// 赋值给外层变量供外面 Response 使用
		createdUserId = userId

		// 如果是第一个注册的用户，系统配置中关闭注册功能
		if userCount == 0 {
			// 写入或更新系统配置项，关闭注册功能（内部需实现 Save 逻辑，支持冷启动 Insert）
			config := &models.SystemConfig{
				Key:    "allow_registration",
				Value:  "false",
				Remark: "是否允许注册",
				Status: 1,
			}
			err = txRepos.SystemConfig.CreateConfig(config)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		response.FailServer(c, "系统繁忙，注册失败: "+err.Error())
		c.Abort()
		return
	}

	response.OkMsg(c, "新用户注册成功，ID为: "+utils.AnyToString(createdUserId))
}
