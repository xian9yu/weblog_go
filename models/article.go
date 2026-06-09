package models

import (
	"errors"
	"time"
	"weblog/dto"

	"gorm.io/gorm"
)

type Article struct {
	ID           uint64         `json:"id" gorm:"not null;primaryKey;unique;comment:文章id"`
	Title        string         `json:"title" gorm:"type:varchar(200);not null;uniqueIndex:idx_title_deleted;comment:文章标题"` // 联合唯一索引
	Content      string         `json:"content" gorm:"type:longtext;not null;comment:文章内容"`
	CategoryName string         `json:"category_name" gorm:"type:varchar(50);not null;comment:分类名称"`
	UserId       uint64         `json:"user_id" gorm:"not null;comment:用户id"`
	Status       int8           `json:"status" gorm:"type:tinyint;not null;comment:是否可以查看，1代表已发布，2代表草稿"`
	CreatedAt    time.Time      `json:"created_at" gorm:"not null;comment:创建时间"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"comment:编辑时间"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;uniqueIndex:idx_title_deleted;precision:3"` // 加上 precision:3，代表将时间戳精确到毫秒（如: 2026-06-09 23:54:46.123）
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
	return article.ID, nil
}

func (repo *ArticleRepository) Update(id uint64, articleData map[string]any) (int64, error) {
	result := repo.db.Model(&Article{}).
		Where("id = ?", id).
		Updates(articleData)
	return result.RowsAffected, result.Error
}

//func (repo *ArticleRepository) Delete(id uint64) int64 {
//	return repo.db.Delete(&Article{}, id).RowsAffected
//}

// BatchDelete 批量删除文章（防越权安全版）
func (repo *ArticleRepository) BatchDelete(ids []uint64, currentUserId uint64) (int64, error) {
	tx := repo.db.Where("id IN ?", ids)

	// 如果不是超级管理员（管理员UID为1）
	// 普通作者只能删除属于自己的文章，防止通过 Postman 拼接 ID 越权删除别人的文章
	if currentUserId != 1 {
		tx = tx.Where("user_id = ?", currentUserId)
	}

	result := tx.Delete(&Article{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil // 返回真正删除了几篇
}

// GetList 纯粹的数据库分页查询逻辑
func (repo *ArticleRepository) GetList(in dto.ArticlePageQueryInput, currentUserId uint64) (total int64, ar []dto.ArticleListResponse, err error) {
	offset := (in.PageNo - 1) * in.PageSize

	// 构建基础查询对象，提前把 LEFT JOIN 挂上去！这样无论是 Count 阶段还是 Scan 阶段，主表和关联表的关系都是完整的
	query := repo.db.Model(&Article{}).
		Joins("LEFT JOIN user ON user.id = article.user_id")

	// 动态拼接 分类 条件
	if in.CategoryName != "" {
		query = query.Where("article.category_name = ?", in.CategoryName)
	}

	// 区分前后台权限划分
	if currentUserId == 0 {
		// 前台游客或公开页面 -> 强制只能看“已发布(1)”
		query = query.Where("article.status = ?", 1)
	} else {
		// 后台管理页面 -> 强制只能看【自己】的资产，但状态不限（草稿、发布都能看）
		query = query.Where("article.user_id = ?", currentUserId)
	}

	// 使用 Session 隔离 Count 句柄，防止污染后续查询
	txCount := query.Session(&gorm.Session{})
	err = txCount.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 如果查出来总数为 0，直接熔断返回，省去一次复杂的 Scan 联表查询 I/O
	if total == 0 {
		return 0, make([]dto.ArticleListResponse, 0), nil
	}

	// 动态排序处理（加上安全白名单防御）
	if in.OrderBy != "" {
		// 限制只能按合法的文章字段排序，且强制统一带上表前缀防歧义
		switch in.OrderBy {
		case "id", "id DESC", "id ASC":
			query = query.Order("article.id DESC")
		case "created_at", "created_at DESC", "created_at ASC":
			query = query.Order("article.created_at DESC")
		case "updated_at", "updated_at DESC", "updated_at ASC":
			query = query.Order("article.updated_at DESC")
		default:
			// 传歪了，降级走默认排序
			query = query.Order("article.created_at DESC")
		}
	} else {
		// 默认按文章创建时间倒序
		query = query.Order("article.created_at DESC")
	}

	// 查列表数据，并自动用 Scan 灌进你的 DTO 结构体
	err = query.Select(`
			article.id, 
			article.title, 
			article.user_id, 
			article.category_name, 
			article.status, 
			article.created_at, 
			article.updated_at,
			user.name AS author_name
		`).
		Limit(in.PageSize).
		Offset(offset).
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

// GetArticleInfosById 通过 id 获取详情
func (repo *ArticleRepository) GetArticleInfosById(articleId uint64) (*Article, error) {
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
