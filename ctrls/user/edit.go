package user

import (
	"github.com/gin-gonic/gin"
	"time"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/encrypt"
	"weblog/utils/response"
)

func Edit(c *gin.Context) {
	id := c.PostForm("id")
	mail := c.PostForm("mail")
	name := c.PostForm("name")
	password := c.PostForm("password")

	//userId := uint64(utils.AnyToInt(id))

	if id == "" {
		response.Error(c, "id 不能为空", nil)
		c.Abort()
		return
	}

	if !utils.VerifyMail(mail) {
		response.Error(c, "邮箱格式不合法", nil)
		c.Abort()
		return
	}

	nameCount := models.CountUserByAny("name", name)
	if nameCount > 0 {
		response.Error(c, "用户名已存在", nil)
		c.Abort()
		return
	}

	mailCount := models.CountUserByAny("mail", mail)
	if mailCount > 0 {
		response.Error(c, "邮箱已存在", nil)
		c.Abort()
		return
	}

	// 获取 token 详情
	userId, _ := c.Get("user_id")
	group, _ := c.Get("group")
	// 只有 (本人/管理员) 才能操作
	if uint64(utils.AnyToInt(id)) == uint64(utils.AnyToInt(userId)) || group == "admin" {
		nowTime := uint64(time.Now().Unix())
		user := models.User{
			Id:          uint64(utils.AnyToInt(id)),
			Name:        name,
			Mail:        mail,
			Password:    encrypt.Md5WithSalt(password, models.PasswordMd5Salt),
			UpdatedTime: nowTime,
		}

		rowsAffected, err := user.Edit()
		if rowsAffected > 0 && err == nil {
			response.Success(c, gin.H{
				"message":       "保存成功",
				"rows_affected": rowsAffected,
			})
			c.Abort()
			return
		}

		response.Error(c, "保存失败, 请稍后再试", err)
		c.Abort()
		return
	}

	response.PageNotFound(c)
}
