package utils

import (
	"math/rand"
	"strconv"
	"time"
	"weblog/utils/encrypt"
)

const ExpireTime int64 = 86400 // token有效期 (单位: s)

var TokenMd5Salt string = v.GetString("securityKey.token")

func NewToken(userName, userId string) string {
	str := userName + "_" + userId + "_" + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.FormatInt(time.Now().Unix(), 10)
	y, m, d := time.Now().Date()
	return encrypt.Md5WithSalt(str, TokenMd5Salt) + AnyToString(y) + AnyToString(m) + AnyToString(d)
}
