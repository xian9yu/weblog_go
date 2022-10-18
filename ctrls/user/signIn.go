package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/encrypt"
	"weblog/utils/response"
)

// 登录
func SignIn(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")

	user, err := models.GetUserDetailsByAny("account", account)
	if err != nil {
		response.Error(c, "登录失败", err)
		c.Abort()
		return
	}

	//if user.Id < 1 {
	//	ctrls.Error(c, "用户不存在", nil)
	//	c.Abort()
	//	return
	//}

	if user.AllowLogin != "y" {
		response.Error(c, "该账号停止使用", nil)
		c.Abort()
		return
	}

	if user.Account == account && user.Password == encrypt.Md5WithSalt(password, models.PasswordMd5Salt) {
		// 事务开始
		var (
			db    = models.MySQL()
			token string
		)
		err = db.Transaction(func(tx *gorm.DB) error {
			// 生成token
			token = utils.NewToken(user.Account, strconv.FormatInt(int64(user.Id), 10))
			auth := models.Auth{
				UserId:    user.Id,
				Token:     token,
				Account:   user.Account,
				Name:      user.Name,
				Mail:      user.Mail,
				Group:     user.Group,
				LoginTime: uint64(time.Now().Unix()),
			}

			// 判断已有登录数据
			if auth.TxCountById(tx) > 0 {
				auth.TxDelete(tx)
			}

			// 登录
			authId, rowsAffected, err := auth.TxAdd(tx)
			if err != nil || rowsAffected < 1 || authId < 1 {
				return err
			}
			return nil
		})

		if err != nil {
			response.Error(c, "登录失败, 请稍后再试", err)
			c.Abort()
			return
		}

		response.Success(c, gin.H{
			"message": "登录成功",
			"token":   token,
		})
		c.Abort()
		return
	}

	response.Error(c, "登录失败, 账户或密码错误", err)
}
