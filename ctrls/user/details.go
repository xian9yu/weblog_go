package user

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

// GetDetailsById 获取详情
func GetDetailsById(c *gin.Context) {
	id := c.Param("id")

	// 获取 token 详情
	userId, _ := c.Get("user_id")
	group, _ := c.Get("group")
	// 只有 (本人/管理员) 才能操作
	if uint64(utils.AnyToInt(id)) == uint64(utils.AnyToInt(userId)) || group == "admin" {
		result, err := models.GetUserDetailsByAny("id", utils.AnyToInt(id))
		if err != nil {
			response.Error(c, "获取用户详情失败", err)
			c.Abort()
			return
		}
		if result.Id < 1 {
			response.Error(c, "用户不存在", nil)
			c.Abort()
			return
		}

		response.Success(c, gin.H{
			"details": result,
		})
		c.Abort()
		return
	}

	response.PageNotFound(c)
}

// GetDetailsByMail 获取详情
func GetDetailsByMail(c *gin.Context) {
	mail := c.Param("mail")

	// 获取 token 详情
	tokenMail, _ := c.Get("mail")
	group, _ := c.Get("group")
	// 只有 (本人/管理员) 才能操作
	if mail == tokenMail || group == "admin" {
		result, err := models.GetUserDetailsByAny("mail", mail)
		if err != nil {
			response.Error(c, "获取用户详情失败", err)
			c.Abort()
			return
		}

		if result.Id < 1 {
			response.Error(c, "用户不存在", nil)
			c.Abort()
			return
		}

		response.Success(c, gin.H{
			"details": result,
		})
		c.Abort()
		return
	}

	response.PageNotFound(c)
}

// 获取详情
func GetDetailsByAccount(c *gin.Context) {
	account := c.Param("account")

	// 获取 token 详情
	tokenAccount, _ := c.Get("mail")
	group, _ := c.Get("group")
	// 只有 (本人/管理员) 才能操作
	if account == tokenAccount || group == "admin" {
		result, err := models.GetUserDetailsByAny("account", account)
		if err != nil {
			response.Error(c, "获取用户详情失败", err)
			c.Abort()
			return
		}

		if result.Id < 1 {
			response.Error(c, "用户不存在", nil)
			c.Abort()
			return
		}

		response.Success(c, gin.H{
			"details": result,
		})
		c.Abort()
		return
	}

	response.PageNotFound(c)
}
