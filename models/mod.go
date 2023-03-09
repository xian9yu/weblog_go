package models

import (
	"weblog/dbConn"
	"weblog/utils"
)

var (
	db = dbConn.MySQL() // 初始化 mysql
	v  = utils.Config() // 初始化配置读取
)

type getDetailsGenerics interface {
	int | string
}
