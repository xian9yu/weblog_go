package dto

import "time"

// ArticleListResponse 用于文章列表展示（不含正文，节省带宽）
type ArticleListResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	UserId       uint64    `json:"user_id"`
	AuthorName   string    `json:"author_name" gorm:"column:author_name"`
	CategoryName string    `json:"category_name"`
	Status       int8      `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ArticleInfoResponse 用于文章详情页（包含正文，但过滤了后台备注等内部字段）
type ArticleInfoResponse struct {
	ID           uint64    `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`       // 文章正文
	Desc         string    `json:"desc"`          // 摘要
	CategoryName string    `json:"category_name"` // 分类名称
	AuthorName   string    `json:"author_name"`   // 返回作者名字
	Status       int8      `json:"status"`        // 状态：0-草稿，1-发布
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
