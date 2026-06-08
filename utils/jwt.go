package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MyClaims 自定义声明结构体，可以把你需要的用户信息塞进 Token 里
type MyClaims struct {
	UserId               string `json:"user_id"`
	UserName             string `json:"user_name"`
	UserEmail            string `json:"user_email"`
	UserGroup            string `json:"user_group"`
	jwt.RegisteredClaims        // 内置的标准声明，包含过期时间等
}

// GenerateToken 签发 JWT Token（支持自定义过期时间）
func GenerateToken(userId, userName, userEmail, userGroup string, expireDuration time.Duration) (string, error) {
	// 创建明文数据
	claims := MyClaims{
		UserId:    userId,
		UserName:  userName,
		UserEmail: userEmail,
		UserGroup: userGroup,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)), // 过期时间
			Issuer:    "weblog_go",                                        // 签发人
		},
	}

	// 用指定的加密算法（常用 HS256）创建 token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用配置文件里的安全密钥进行签名，生成最终的字符串 Token
	secretKey := []byte(GetConfig().SecurityKey.Token)
	return token.SignedString(secretKey)
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenString string) (*MyClaims, error) {
	secretKey := []byte(GetConfig().SecurityKey.Token)

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 校验 Token 是否有效并转换成我们的自定义结构体
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
