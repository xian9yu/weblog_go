package utils

import (
	"regexp"
)

// 1. 在包级别预编译正则表达式，整个程序运行期间只编译一次
// 使用 MustCompile，如果在程序启动时正则写错了，会直接 panic 报错，保证线上安全
var mailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// CheckEmail 验证邮箱格式是否合法（高性能、无内存抖动版）
func CheckEmail(mail string) bool {
	// 如果传入空字符串，直接返回 false，避免走正则匹配
	if mail == "" {
		return false
	}

	// 2. 直接复用编译好的 mailRegexp 对象，极速匹配
	return mailRegexp.MatchString(mail)
}
