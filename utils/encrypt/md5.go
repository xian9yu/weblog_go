package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5WithSalt 生成md5加盐字串
func Md5WithSalt(value, salt string) string {
	m := md5.New()
	m.Write([]byte(value + salt))
	return hex.EncodeToString(m.Sum(nil))
}

// Md5 encode with md5 string
func Md5Encode(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}
