package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"
	"weblog/utils"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth 普通鉴权
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头（Header）中获取 token
		token := c.GetHeader("Authorization")
		if token == "" {
			response.FailUnauthorized(c, "请求未携带Token，请先登录")
			c.Abort() // 必须 Abort，阻止执行后面的业务 Controller
			return
		}

		// 检查 Token 格式是否为标准的 "Bearer <token>"
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.FailUnauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		// 解析并验证 Token 的合法性与有效性
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			response.FailUnauthorized(c, "登录令牌无效或已过期，请重新登录")
			c.Abort()
			return
		}

		// 计算一下当前请求 Token 的 MD5
		tokenStr := parts[1]
		hasher := md5.New()
		hasher.Write([]byte(tokenStr))
		tokenMd5 := hex.EncodeToString(hasher.Sum(nil))
		//  拦截token：如果黑名单里有它，说明点了退出登录！
		if utils.Blacklist.Contains(tokenMd5) {
			response.FailUnauthorized(c, "登录已失效，请重新登录")
			c.Abort() // 必须 Abort，拒绝他继续访问后面的 Controller
			return
		}

		// 从 claims 中拿到这个 Token 真正的过期时间
		// claims.ExpiresAt 是 jwt 库内置的 *jwt.NumericDate 类型，用 .Time 转换回 Go 的 time.Time
		expireTime := claims.ExpiresAt.Time

		// 计算【当前时间】距离【过期时间】还剩多少秒
		remainingTime := time.Until(expireTime)

		// 如果剩余时间大于 0（没过期），且小于等于 10 分钟 (10 * time.Minute)
		if remainingTime > 0 && remainingTime <= 10*time.Minute {
			// 触发自动续期：重新签发一张满血的 8 小时有效期的 Token
			newToken, err := utils.GenerateToken(claims.UserId, claims.UserName, claims.UserEmail, claims.UserGroup, 8*time.Hour)
			if err == nil {
				// 将新 Token 悄悄塞进 Response 的 Header 里返回给前端
				// 格式同样保持标准的 Bearer 格式
				c.Header("X-Refresh-Token", "Bearer "+newToken)

				// ⚠️ 跨域安全提示：如果是前后端分离的项目，必须暴露这个自定义 Header，否则前端 JS 无法读取
				c.Header("Access-Control-Expose-Headers", "X-Refresh-Token")
			}
		}

		// 将解析出来的用户信息塞进 Gin 的上下文 Context 里
		c.Set("user_id", utils.AnyToUint64(claims.UserId))
		c.Set("name", claims.UserName)
		c.Set("email", claims.UserEmail)
		c.Set("group", claims.UserGroup)

		c.Next() // 放行，去执行具体的业务路由
	}
}
