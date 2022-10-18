package utils

import (
	"errors"
	"os"
)

// 判断文件夹是否存在
func hasDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true, nil
		}
		return false, err
	}
	return true, err
}

// 创建文件夹
func CreateDir(path string) error {
	exist, err := hasDir(path)
	if err != nil {
		return errors.New("获取文件夹异常: " + err.Error())
	}
	if !exist {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return errors.New("创建目录异常: " + err.Error())
		} else {
			return nil
		}
	}

	return errors.New("文件夹已存在")
}
