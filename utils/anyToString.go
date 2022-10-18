package utils

import (
	"encoding/json"
	"strconv"
)

// AnyToString any 转 string
func AnyToString(value any) (stringValue string) {
	if value == nil {
		return stringValue
	}

	switch v := value.(type) {
	case float64:
		stringValue = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		stringValue = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case int:
		stringValue = strconv.Itoa(v)
	case uint:
		stringValue = strconv.Itoa(int(v))
	case int8:
		stringValue = strconv.Itoa(int(v))
	case uint8:
		stringValue = strconv.Itoa(int(v))
	case int16:
		stringValue = strconv.Itoa(int(v))
	case uint16:
		stringValue = strconv.Itoa(int(v))
	case int32:
		stringValue = strconv.Itoa(int(v))
	case uint32:
		stringValue = strconv.Itoa(int(v))
	case int64:
		stringValue = strconv.FormatInt(v, 10)
	case uint64:
		stringValue = strconv.FormatUint(v, 10)
	case string:
		stringValue = v
	case []byte:
		stringValue = string(v)
	default:
		newValue, _ := json.Marshal(value)
		stringValue = string(newValue)
	}

	return stringValue
}
