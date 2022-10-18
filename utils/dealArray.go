package utils

import (
	"fmt"
	"sort"
	"strings"
)

// IsContainString 判断元素是否在数组中
func IsContainString(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// IsContainInt64 判断元素是否在数组中
func IsContainInt64(items []int64, item int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// IdsToString 将查询出来的数组转成字符串格式 example:[1,2,3] => "1,2,3"
func IdsToString(ids []int64) string {
	return strings.Replace(strings.Trim(fmt.Sprint(ids), "[]"), " ", ",", -1)
}

// StringIdsToString 去除id为0的元素
func StringIdsToString(ids []string) string {
	var newIds []string
	for _, v := range ids {
		if v != "0" {
			newIds = append(newIds, v)
		}
	}
	return strings.Replace(strings.Trim(fmt.Sprint(newIds), "[]"), " ", ",", -1)
}

// SliceRemoveDuplicates 去除数组中相同元素,这种发放适用于string,int,float等切片，会对切片中的元素进行排序
func nn(slice []string) []string {
	sort.Strings(slice)
	i := 0
	var j int
	for {
		if i >= len(slice)-1 {
			break
		}

		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {
		}
		slice = append(slice[:i+1], slice[j:]...)
		i++
	}
	return slice
}

// RemoveDuplicate 去除数组中相同元素
func RemoveDuplicate(list *[]int64) []int64 {
	var x []int64
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}
