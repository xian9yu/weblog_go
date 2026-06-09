package utils

import (
	"regexp"
)

// 1. 在包加载时预编译，整个程序运行期间只编译一次
// 同步国内最新手机号段：包含 13/14/15/16/17/18/19 下属的所有合法虚拟和实体运营商号段
var mobileRegexp = regexp.MustCompile(`^1[3-9]\d{9}$`)

// CheckMobile 验证手机号格式是否合法（高性能、最新号段版）
func CheckMobile(mobileNum string) bool {
	// 前置防御
	if mobileNum == "" {
		return false
	}

	// 2. 直接复用编译好的 mobileRegexp 对象，极速匹配
	return mobileRegexp.MatchString(mobileNum)
}
