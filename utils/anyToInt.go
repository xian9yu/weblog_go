package utils

import (
	"strconv"
)

// AnyToInt any 转 int
func AnyToInt(value any) (intValue int) {
	switch v := value.(type) {
	case uint:
		intValue = int(v)
	case int8:
		intValue = int(v)
	case uint8:
		intValue = int(v)
	case int16:
		intValue = int(v)
	case uint16:
		intValue = int(v)
	case int32:
		intValue = int(v)
	case uint32:
		intValue = int(v)
	case int64:
		intValue = int(v)
	case uint64:
		intValue = int(v)
	case float32:
		intValue = int(v)
	case float64:
		intValue = int(v)
	case string:
		intValue, _ = strconv.Atoi(v) // string 不能是 float 类型
	default:
		intValue = value.(int)
	}
	return intValue
}
