package models

import "gorm.io/gorm"

type Upload struct {
	Id          uint64 `json:"id" gorm:"primaryKey;comment:id"`
	UserId      uint64 `json:"user_id" gorm:"not null;comment:上传用户id"`
	Name        string `json:"name" gorm:"not null;comment:名称"`
	FullName    string `json:"full_name" gorm:"not null;comment:完整名称(带文件格式)"`
	Path        string `json:"prefix" gorm:"not null;comment:本地存储path"`
	FullUrl     string `json:"full_url" gorm:"not null;comment:完整下载url"`
	Size        uint64 `json:"size" gorm:"not null;comment:大小"`
	Suffix      string `json:"suffix" gorm:"not null;comment:文件格式"`
	Category    string `json:"category" gorm:"not null;comment:文件分类"`
	Remark      string `json:"remark" gorm:"comment:备注"`
	CreatedTime uint64 `json:"created_time" gorm:"not null;comment:创建时间"`
}

func (upload Upload) Add() (id uint64, rowsAffected int64, err error) {
	result := db.Create(&upload)
	return upload.Id, result.RowsAffected, result.Error
}

func (Upload) Delete(id uint64) int64 {
	return db.Delete(&Upload{}, id).RowsAffected
}

func (Upload) List() (total int64, upload []Upload, err error) {
	db.Model(&Upload{}).Count(&total)
	err = db.Find(&upload).Error
	return total, upload, err
}

func GetUploadDetailsByAny[T getDetailsGenerics](key string, value T) (upload Upload, err error) {
	err = db.Where(key+" = ?", value).Find(&upload).Error
	return upload, err
}

func CountUploadByAny[T getDetailsGenerics](key string, value T) (count int64) {
	db.Model(&Upload{}).Where(key+" = ?", value).Count(&count)
	return count
}

// 事务

func (upload Upload) TxGetDetailsById(tx *gorm.DB) (u Upload, err error) {
	err = tx.Where("id = ?", upload.Id).Find(&u).Error
	return u, err
}

func (upload Upload) TxDelete(tx *gorm.DB) int64 {
	return tx.Delete(&Upload{}, upload.Id).RowsAffected
}
