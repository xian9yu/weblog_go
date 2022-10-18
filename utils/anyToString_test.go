package utils

import (
	"fmt"
	"testing"
)

func TestAnyToString(t *testing.T) {
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
	var m string = "stringValue"
	var n []byte = []byte("safeguards")

	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(a), AnyToString(a)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(b), AnyToString(b)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(c), AnyToString(c)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(d), AnyToString(d)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(e), AnyToString(e)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(f), AnyToString(f)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(g), AnyToString(g)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(h), AnyToString(h)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(i), AnyToString(i)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(j), AnyToString(j)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(k), AnyToString(k)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(l), AnyToString(l)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(m), AnyToString(m)))
	t.Log(fmt.Sprintf("type = %T , value = %s", AnyToString(n), AnyToString(n)))

}
