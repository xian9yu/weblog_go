package middleware

//import (
//	"time"
//	"weblog/models"
//	"weblog/utils"
//	"weblog/utils/response"
//
//	"github.com/gin-gonic/gin"
//)
//
//// AdminAuth 管理员专用鉴权
//func AdminAuth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token := c.Request.Header.Get("token")
//
//		if token == "" {
//			response.FailUnauthorized(c, "请求未携带token")
//	//			return
//		}
//
//		// 判断 token 是否存在
//		auth := models.Auth{
//			Token: token,
//		}
//		authDetails, err := auth.GetInfoByToken()
//		if err != nil {
//			response.FailUnauthorized(c, "token失效，请重新登录")
//	//			return
//		}
//
//		// 判断 token 过期时间
//		if uint64(time.Now().Unix())-authDetails.LoginTime > uint64(utils.ExpireTime) {
//			response.FailUnauthorized(c, "token已过期，请重新登录")
//	//			return
//		}
//
//		// 限制只允许管理员操作
//		if authDetails.Group != "admin" {
//			response.FailUnauthorized(c)
//	//			return
//		}
//
//		// 判断 token有效时间小于配置中设定时间的三分之一则更新过期时间
//		if uint64(time.Now().Unix())-authDetails.LoginTime < uint64(3/utils.ExpireTime) {
//			// 更新 token过期时间
//			auth.LoginTime = uint64(time.Now().Unix())
//			_, _ = auth.Update()
//		}
//
//		c.Set("token", token)
//		c.Set("user_id", authDetails.UserId)
//		c.Set("username", authDetails.Name)
//		c.Set("mail", authDetails.Email)
//		c.Set("group", authDetails.Group)
//		c.Set("login_time", authDetails.LoginTime)
//	}
//}
