package models

import (
	"gorm.io/gorm"
	"weblog/utils"
)

type Article struct {
	Id          uint64 `json:"id" gorm:"not null;primaryKey;unique;comment:文章id"`
	Title       string `json:"title" gorm:"not null;unique;comment:文章标题"`
	Content     string `json:"content" gorm:"not null;comment:文章内容"`
	AuthorId    uint64 `json:"author_id" gorm:"not null;comment:作者"`
	AllowView   string `json:"allow_view" gorm:"not null;comment:是否可以查看 n隐藏 y显示"`
	CreatedTime uint64 `json:"created_time" gorm:"not null;comment:创建时间"`
	UpdatedTime uint64 `json:"updated_time" gorm:"not null;comment:编辑时间"`
}

func (article Article) Add() (uint64, int64, error) {
	result := db.Create(&article)
	return article.Id, result.RowsAffected, result.Error
}

func (article Article) Edit() (int64, error) {
	result := db.Model(&Article{}).Where("id = ?", article.Id).Updates(article)
	return result.RowsAffected, result.Error
}

func (Article) Delete(id uint64) int64 {
	return db.Delete(&Article{}, id).RowsAffected
}

// 分类列表
func (Article) GetCategoryList(articleIds []int, orderBy string, page utils.Pagination) (total int64, article []Article, err error) {
	var (
		size   = page.LimitSize
		offset = page.Offset()
	)

	if orderBy == "" {
		orderBy = "created_time DESC"
	}

	db.Model(&Article{}).Where("allow_view = ? and id in ?", "y", articleIds).Count(&total)
	err = db.Where(Article{}).Where("allow_view = ? and id in ?", "y", articleIds).Limit(size).Offset(offset).Order(orderBy).Find(&article).Error

	return total, article, err
}

// 全列表
func (Article) GetList(orderBy string, page utils.Pagination) (total int64, article []Article, err error) {
	var (
		size   = page.LimitSize
		offset = page.Offset()
	)

	if orderBy == "" {
		orderBy = "created_time DESC"
	}

	db.Model(&Article{}).Where("allow_view = ?", "y").Count(&total)
	err = db.Where(Article{}).Where("allow_view = ?", "y").Limit(size).Offset(offset).Order(orderBy).Find(&article).Error

	return total, article, err
}

// CountArticleByAny 获取有状态的count
func CountArticleByAny[T getDetailsGenerics](key string, value T) (count int64) {
	db.Model(&Article{}).Where("allow_view = ? and "+key+" = ?", "y", value).Count(&count)
	return count
}

// CountArticleByAny 获取所有数量(无限制)
//func CountArticleByAny[T getDetailsGenerics](str string, value T) (count int64) {
//	db.Model(&Article{}).Where(str+" = ?", value).Count(&count)
//	return count
//}

// 获取详情
func GetArticleDetailsByAny[T getDetailsGenerics](key string, value T) (a Article, err error) {
	var article Article
	err = db.Where("allow_view = ? and "+key+" = ?", "y", value).Find(&article).Error
	return article, err
}

// 事务

func (article Article) TxAdd(tx *gorm.DB) (uint64, int64, error) {
	result := tx.Create(&article)
	return article.Id, result.RowsAffected, result.Error
}

func (article Article) TxEdit(tx *gorm.DB) (int64, error) {
	result := tx.Model(&Article{}).Where("id = ?", article.Id).Updates(article)
	return result.RowsAffected, result.Error
}
