package models

import (
	"errors"
	"weblog/dto"

	"gorm.io/gorm"
)

type Article struct {
	Id           uint64 `json:"id" gorm:"not null;primaryKey;unique;comment:文章id"`
	Title        string `json:"title" gorm:"not null;unique;comment:文章标题"`
	Content      string `json:"content" gorm:"not null;comment:文章内容"`
	CategoryName string `json:"category_name" gorm:"not null;comment:分类名称"`
	UserId       uint64 `json:"user_id" gorm:"not null;comment:用户id"`
	Status       int8   `json:"status" gorm:"not null;comment:是否可以查看 0代表草稿，1代表已发布"`
	CreatedTime  uint64 `json:"created_time" gorm:"not null;comment:创建时间"`
	UpdatedTime  uint64 `json:"updated_time" gorm:"not null;comment:编辑时间"`
}

// ArticleRepository 封装所有对 article 表的底层 SQL 操作
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 初始化时注入 DB 实例
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// WithTx 产生一个带有事务连接的新 Repo 实例
func (repo *ArticleRepository) WithTx(tx *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: tx} // 将原本的 db 替换为事务 tx
}

func (repo *ArticleRepository) Add(article Article) (uint64, error) {
	result := repo.db.Create(&article)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected < 1 {
		return 0, gorm.ErrRecordNotFound
	}
	return article.Id, nil
}

func (repo *ArticleRepository) Edit(id uint64, articleData map[string]any) (int64, error) {
	result := repo.db.Model(&Article{}).
		Where("id = ?", id).
		Updates(articleData)
	return result.RowsAffected, result.Error
}

func (repo *ArticleRepository) Delete(id uint64) int64 {
	return repo.db.Delete(&Article{}, id).RowsAffected
}

// GetList 纯粹的数据库分页查询逻辑
func (repo *ArticleRepository) GetList(articleInput dto.ArticlePageQueryInput, currentUserId uint64) (total int64, ar []dto.ArticleListResponse, err error) {
	offset := (articleInput.Page - 1) * articleInput.PageSize

	query := repo.db.Model(&Article{}) // 构建基础查询对象

	// 动态拼接 分类 条件
	if articleInput.CategoryName != "" {
		query = query.Where("article.category_name = ?", articleInput.CategoryName)
	}

	if currentUserId == 0 {
		// 前台游客或公开页面 -> 强制只能看“已发布(1)”
		query = query.Where("article.status = ?", 1)
	} else {
		// 后台管理页面 -> 强制只能看【自己】的资产，但状态不限（草稿、发布都能看）
		query = query.Where("article.user_id = ?", currentUserId)
	}

	// 动态排序处理（带表前缀防歧义，并提供默认值）
	if articleInput.OrderBy != "" {
		query = query.Order(articleInput.OrderBy)
	} else {
		// 默认按文章创建时间倒序
		query = query.Order("article.created_time DESC")
	}

	err = query.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}
	// 查列表数据，并自动用 Scan 灌进你的 DTO 结构体
	err = query.Select(`
			article.id, 
			article.title, 
			article.user_id, 
			article.category_name, 
			article.status, 
			article.created_time, 
			article.updated_time,
			user.name AS author_name
		`).
		Joins("LEFT JOIN user ON user.id = article.user_id").
		Limit(int(articleInput.PageSize)).
		Offset(int(offset)).
		Scan(&ar).Error

	return total, ar, err
}

// CountArticleByTitle 获取title的count
func (repo *ArticleRepository) CountArticleByTitle(title string) (count int64) {
	repo.db.Model(&Article{}).Where("title = ?", title).Count(&count)
	return count
}

// ArticleTitleExists 判断 title 是否重复（排除指定 ID）
func (repo *ArticleRepository) ArticleTitleExists(title string, excludeId uint64) (bool, error) {
	var id uint64
	// 只查询 ID 字段，并且只查 1 条（Limit 1），速度快到飞起
	err := repo.db.Model(&Article{}).Select("id").Where("title = ? AND id != ?", title, excludeId).Limit(1).Pluck("id", &id).Error
	if err != nil {
		// 如果是单纯的“没查到记录”，在 GORM v1/v2 中某些查询会触发 ErrRecordNotFound
		// 但使用 Pluck 查空切片时通常 err 为 nil。为了绝对安全，做个兼容判断：
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		// 如果是断网、语法错误等真正的数据库异常，抛给上层 Controller
		return false, err
	}

	// 如果 id > 0 说明记录存在
	return id > 0, nil

}

// CountArticleByAny 获取所有数量(无限制)
//func CountArticleByAny[T getDetailsGenerics](str string, value T) (count int64) {
//	database.Model(&Article{}).Where(str+" = ?", value).Count(&count)
//	return count
//}

// GetArticleDetailsById 通过 id 获取详情
func (repo *ArticleRepository) GetArticleDetailsById(articleId uint64) (*Article, error) {
	var article Article
	err := repo.db.Where("allow_view = ? and id = ?", "y", articleId).Find(&article).Error
	return &article, err
}

// GetArticleUserIdById 仅仅查询文章的作者 ID
func (repo *ArticleRepository) GetArticleUserIdById(id uint64) (uint64, error) {
	var userId uint64
	err := repo.db.Model(&Article{}).
		Select("user_id").
		Where("id = ?", id).
		Row().
		Scan(&userId) // 直接把这一列的值扫描到变量里

	if err != nil {
		return 0, err
	}
	return userId, nil
}

// 事务

//func (article Article) TxAdd(tx *gorm.DB) (uint64, int64, error) {
//	result := tx.Create(&article)
//	return article.Id, result.RowsAffected, result.Error
//}

//func (article Article) TxEdit(tx *gorm.DB) (int64, error) {
//	result := tx.Model(&Article{}).Where("id = ?", article.Id).Updates(article)
//	return result.RowsAffected, result.Error
//}
