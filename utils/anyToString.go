package utils

import (
	"fmt"
	"strconv"
)

// AnyToString 将任意基本类型安全、高性能地转换为 string
func AnyToString(value any) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint8: // ⚠️ 注意：Go 里的 byte 就是 uint8，但如果传入 []byte 数组需要单独处理
		return strconv.FormatUint(uint64(v), 10)
	case bool:
		return strconv.FormatBool(v)
	case float64:
		// 'f' 代表不使用科学计数法，-1 代表保持原样精度，避免 1000000 变成 1e+06
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case []byte:
		return string(v)
	default:
		// 最后的底层兜底，处理结构体、Map等复杂类型，虽然慢点但胜在能跑
		return fmt.Sprintf("%v", v)
	}
}
