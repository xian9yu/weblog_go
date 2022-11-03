package models

import (
	"gorm.io/gorm"
	"weblog/utils"
)

var PasswordMd5Salt string = v.GetString("securityKey.password")

type User struct {
	Id          uint64 `json:"id" gorm:"size:12;primaryKey;unique;notnull;comment:用户id"`
	Account     string `json:"account" gorm:"size:60;comment:账号"`
	Name        string `json:"name" gorm:"size:60;comment:用户名"`
	Mail        string `json:"mail" gorm:"size:60;notnull;comment:邮箱"`
	Password    string `json:"-"  gorm:"size:33;notnull;comment:登录密码"`
	Group       string `json:"group" gorm:"size:9;notnull;comment:用户分组(管理员;用户)"`
	AllowLogin  string `json:"allow_login" gorm:"size:1;notnull;comment:允许登录(y启用;n关闭)"`
	CreatedTime uint64 `json:"created_time" gorm:"autoCreateTime;notnull comment:user创建时间"`
	UpdatedTime uint64 `json:"updated_time" gorm:"autoUpdateTime;comment:上一次修改信息时间"`
}

// Login 用户登录
func (user User) Login() bool {
	var count int64
	db.Model(&User{}).Where("account = ? AND password = ?", &user.Account, &user.Password).Count(&count)
	return count > 0
}

// 注册
func (user User) Add() (userId uint64, rowsAffected int64, err error) {
	result := db.Create(&user)
	return user.Id, result.RowsAffected, result.Error
}

func (user User) Edit() (rowsAffected int64, err error) {
	result := db.Model(&User{}).Where("id = ?", user.Id).Updates(user)
	return result.RowsAffected, result.Error
}

func (User) Delete(id uint64) int64 {
	return db.Delete(&User{}, id).RowsAffected
}

// 查询数量
func CountUserByAny[T getDetailsGenerics](key string, value T) (count int64) {
	db.Model(&User{}).Where(key+" = ?", value).Count(&count)
	return count
}

// GetUserDetailsByAny 获取详情
func GetUserDetailsByAny[T getDetailsGenerics](key string, value T) (user User, err error) {
	err = db.Where(key+" = ?", value).Find(&user).Error
	return user, err
}

// 查列表
func (User) GetList(group, orderBy string, page utils.Pagination) (total int64, user []User, err error) {
	var (
		size   = page.PageSize
		offset = page.Offset()
	)

	if orderBy == "" {
		orderBy = "created_time ASC"
	}

	queryDB := db.Model(&User{})

	switch group {
	case "admin":
		db.Model(&User{}).Where("`group` = ?", group).Count(&total)
		queryDB = db.Where("`group` = ?", group)
	case "user":
		db.Model(&User{}).Where("`group` = ?", group).Count(&total)
		queryDB = db.Where("`group` = ?", group)
	default:
		db.Model(&User{}).Count(&total)
	}

	err = queryDB.Limit(size).Offset(offset).Order(orderBy).Find(&user).Error

	return total, user, err
}

// 事务

func (user User) TxAdd(tx *gorm.DB) (userId uint64, rowsAffected int64, err error) {
	result := tx.Create(&user)
	return user.Id, result.RowsAffected, result.Error
}

func (user User) TxEdit(tx *gorm.DB, id uint64) (rowsAffected int64, err error) {
	result := tx.Model(&User{}).Where("id = ?", id).Updates(user)
	return result.RowsAffected, result.Error
}
