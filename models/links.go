package models

import (
	"gorm.io/gorm"
)

type Links struct {
	ArticleId   uint64 `json:"article_id" gorm:"not null;comment:文章id"`
	CategoryId  uint64 `json:"category_id" gorm:"not null;comment:分类id"`
	UpdatedTime uint64 `json:"updated_time" gorm:"not null;comment:编辑时间"`
}

func (Links) Add(links []Links) (rowsAffected int64, err error) {
	result := db.Create(&links)
	return result.RowsAffected, result.Error
}

// 精准删除
func (Links) Delete(articleId, categoryId uint64) (rowsAffected int64) {
	return db.Where("article_id = ? and category_id = ?", articleId, categoryId).Delete(&Links{}).RowsAffected
}

// 通过分类id删除
func (links Links) DeleteByCategoryId() (rowsAffected int64) {
	return db.Where("category_id = ?", links.CategoryId).Delete(&Links{}).RowsAffected
}

// 通过文章id删除
func (links Links) DeleteByArticleId() (rowsAffected int64) {
	return db.Where("article_id = ?", links.ArticleId).Delete(&Links{}).RowsAffected
}

func (links Links) CountLinksById() (count int64) {
	db.Model(&Links{}).Where("article_id = ? and category_id = ?", links.ArticleId, links.CategoryId).Count(&count)
	return count
}

// GetList 列表
func (Links) GetList(categoryIds []int) (links []Links, err error) {
	err = db.Model(&Links{}).Where("category_id in ?", categoryIds).Find(&links).Error
	return links, err
}

// 事务

func (Links) TxAdd(tx *gorm.DB, links []Links) (rowsAffected int64, err error) {
	result := tx.Create(&links)
	return result.RowsAffected, result.Error
}

// 通过 article id  查 分类列表
func (links Links) TxGetCategoryIdByArticleId(tx *gorm.DB) (l []Links, err error) {
	err = tx.Model(&Links{}).Where("article_id = ?", links.ArticleId).Find(&l).Error
	return l, err
}

// 通过分类id删除
func (links Links) TxDeleteByCategoryId(tx *gorm.DB) error {
	return tx.Where("category_id = ?", links.CategoryId).Delete(&Links{}).Error
}

// 通过文章id删除
func (links Links) TxDeleteByArticleId(tx *gorm.DB) error {
	return tx.Where("article_id = ?", links.ArticleId).Delete(&Links{}).Error
}
