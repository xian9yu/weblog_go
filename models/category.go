package models

import "gorm.io/gorm"

type Category struct {
	Id          uint64 `json:"id" gorm:"primaryKey;comment:id"`
	Name        string `json:"name" gorm:"not null;comment:名称"`
	Type        string `json:"type" gorm:"not null;comment:分类"`
	Remark      string `json:"remark" gorm:"comment:备注"`
	CreatedTime uint64 `json:"created_time" gorm:"not null;comment:创建时间"`
	UpdatedTime uint64 `json:"updated_time" gorm:"not null;comment:编辑时间"`
}

func (category Category) Add() (categoryId uint64, rowsAffected int64, err error) {
	result := db.Create(&category)
	return category.Id, result.RowsAffected, result.Error
}

func (category Category) Edit() (int64, error) {
	result := db.Model(&Category{}).Where("id = ?", category.Id).Updates(category)
	return result.RowsAffected, result.Error
}

func (Category) Delete(id uint64) int64 {
	return db.Delete(&Category{}, id).RowsAffected
}

// List 列表
func (Category) List() (total int64, category []Category, err error) {
	db.Model(&Category{}).Count(&total)
	err = db.Select("id, name, type, remark").Find(&category).Error
	return total, category, err
}

func CountCategoryByAny[T getDetailsGenerics](key string, value T) (count int64) {
	db.Model(&Category{}).Where(key+" = ?", value).Count(&count)
	return count
}

// 获取详情
func GetCategoryDetailsByAny[T getDetailsGenerics](key string, value T) (category Category, err error) {
	err = db.Where(key+"= ?", value).Find(&category).Error
	return category, err
}

// 事务

//func (Category) IdList(tx *gorm.DB) (total int64, category []Category, err error) {
//	tx.Model(&Category{}).Where("`type` = ?", "category").Count(&total)
//	err = tx.Select("id").Where("`type` = ?", "category").Find(&category).Error
//	return total, category, err
//}

func (Category) TxDelete(tx *gorm.DB, id uint64) error {
	return tx.Delete(&Category{}, id).Error
}
