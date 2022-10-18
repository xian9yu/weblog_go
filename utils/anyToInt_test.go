package utils

import (
	"fmt"
	"testing"
)

func TestAnyToInt(t *testing.T) {
	var a int = 6
	var b uint = 3
	var c int8 = 3
	var d uint8 = 3
	var e int16 = 3
	var f uint16 = 3
	var g int32 = 3
	var h uint32 = 3
	var i int64 = 3
	var j uint64 = 3
	var k float32 = 1121.453
	var l float64 = 111.31
	var m string = "111" // 不能为 float

	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(a), AnyToInt(a)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(b), AnyToInt(b)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(c), AnyToInt(c)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(d), AnyToInt(d)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(e), AnyToInt(e)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(f), AnyToInt(f)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(g), AnyToInt(g)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(h), AnyToInt(h)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(i), AnyToInt(i)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(j), AnyToInt(j)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(k), AnyToInt(k)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(l), AnyToInt(l)))
	t.Log(fmt.Sprintf("type = %T , value = %d", AnyToInt(m), AnyToInt(m)))

}
