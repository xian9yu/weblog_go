package utils

import (
	"time"
)

func CurrentTime() int64 {
	return time.Now().UnixNano() / 1000000
}

func FormatTime() string {
	return time.Now().Format("2006.01.02 15:04:05")
}

func FormatTimeYYMMDD() string {
	return time.Now().Format("2006.01.02")
}

func FormatTimeYYMM() string {
	return time.Now().Format("2006.01")
}

func FormatTimeMMDD() string {
	return time.Now().Format("0102")
}

//println("ss")
//fmt.Printf("时间戳（秒）：%v;\n", time.Now().Unix())
//fmt.Printf("时间戳（纳秒）：%v;\n",time.Now().UnixNano())
//fmt.Printf("时间戳（毫秒）：%v;\n",time.Now().UnixNano() / 1e6)
//fmt.Printf("时间戳（纳秒转换为秒）：%v;\n",time.Now().UnixNano() / 1e9)
//
//t := time.Now()
//fmt.Printf("%04d.%02d.%2d\n", t.Year(), t.Month(), t.Day())
