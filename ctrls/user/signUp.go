package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/encrypt"
	"weblog/utils/response"
)

// 注册
func SignUp(c *gin.Context) {
	account := c.PostForm("account")
	name := c.PostForm("name")
	password := c.PostForm("password")
	mail := c.PostForm("mail")

	// 参数是否为空
	if account == "" {
		response.Error(c, "account 不能为空", nil)
		c.Abort()
		return
	}
	if password == "" {
		response.Error(c, "密码不能为空", nil)
		c.Abort()
		return
	}

	// 校验邮箱格式
	if !utils.VerifyMail(mail) {
		response.Error(c, "邮箱格式不合法", nil)
		c.Abort()
		return
	}

	// 判断是否允许注册
	setting, er := models.GetSettingDetailsByAny("name", "allow_signup")
	if er != nil {
		response.Error(c, "注册失败, 请重试", er)
		c.Abort()
		return
	}
	if setting.Value != "y" {
		response.Error(c, "未开启注册", nil)
		c.Abort()
		return
	}

	// 判断 用户名|邮箱 存在
	mailCount := models.CountUserByAny("mail", mail)
	if mailCount > 0 {
		response.Error(c, "邮箱已存在", nil)
		c.Abort()
		return
	}
	accountCount := models.CountUserByAny("account", account)
	if accountCount > 0 {
		response.Error(c, "account 已存在", nil)
		c.Abort()
		return
	}

	nowTime := uint64(time.Now().Unix())
	user := &models.User{
		Account:     account,
		Name:        name,
		Mail:        mail,
		Password:    encrypt.Md5WithSalt(password, models.PasswordMd5Salt),
		Group:       "user",
		AllowLogin:  "y",
		CreatedTime: nowTime,
		UpdatedTime: nowTime,
	}

	// 事务开始
	var (
		db  = models.MySQL()
		uId uint64
	)
	err := db.Transaction(func(tx *gorm.DB) error {
		// 新增用户
		userId, rowsAffected, err := user.TxAdd(tx)
		if rowsAffected < 1 || err != nil {
			return err
		}
		uId = userId

		//如果新增用户id为1, 则设置为管理员账户
		if userId == 1 {
			u := models.User{
				Group:       "admin",
				UpdatedTime: uint64(time.Now().Unix()),
			}
			affected, err := u.TxEdit(tx, userId)
			if affected < 1 || err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		response.Error(c, "用户注册失败, 请稍后再试", err)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"message": "注册成功",
		"user_id": uId,
	})
}
