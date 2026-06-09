package dto

import "time"

// ConfigDetailOutput 配置详情的返回数据（出参）
type ConfigDetailOutput struct {
	Id        uint64    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConfigListOutput 配置项列表的统一分页返回结构
type ConfigListOutput struct {
	List  []ConfigDetailOutput `json:"list"`  // 配置项数组
	Total int64                `json:"total"` // 总条数，方便前端渲染分页组件
}
