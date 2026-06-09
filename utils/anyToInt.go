package utils

import (
	"strconv"
)

// AnyToInt64 将任意基本类型安全、高性能地转换为 int64
func AnyToInt64(value any) int64 {
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		// 优先使用高性能的 ParseInt
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
		// 如果是带小数点的字符串（如 "12.34"），尝试当成浮点数解析再转整
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(f)
		}
		return 0
	case float64:
		return int64(v) // 核心：解决 JSON 解析或高精度计算后转整的痛点
	case float32:
		return int64(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	case uint:
		return int64(v)
	case uint64:
		return int64(v)
	case uint32:
		return int64(v)
	case uint16:
		return int64(v)
	case uint8:
		return int64(v)
	case []byte: // 兼容从数据库（如文本类型）或缓存中读出的字节流
		if i, err := strconv.ParseInt(string(v), 10, 64); err == nil {
			return i
		}
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func AnyToInt(value any) int {
	return int(AnyToInt64(value))
}
