package dto

// ArticleCreateInput 创建文章时的入参校验
type ArticleCreateInput struct {
	Title        string `json:"title" binding:"required,min=3,max=100"`  // min=3,max=100：标题长度必须在 3-100 个字符之间
	Content      string `json:"content" binding:"required"`              // required：文章内容不能为空
	CategoryName string `json:"category_name" binding:"required,max=30"` // max=30：分类名称不能太长，防止长文本撑破数据库或前端排版
	Status       int8   `json:"status" binding:"oneof=0 1"`              // oneof=0 1：限定状态值，0 代表草稿，1 代表已发布。传 2 或 -1 直接报错拦截
}

// ArticleUpdateInput 更新文章时的入参校验
type ArticleUpdateInput struct {
	ID uint64 `json:"id" binding:"required,gt=0"` // ID 是必填的，用来确定更新哪一条
	// 更新时由于支持局部更新，很多字段不需要必填(required)，但一旦传了，就必须满足长度和范围限制
	Title        string `json:"title" binding:"omitempty,min=3,max=100"`
	Content      string `json:"content" binding:"omitempty"`
	CategoryName string `json:"category_name" binding:"omitempty,max=30"`
	Status       int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

// ArticlePageQueryInput 文章分页列表查询参数校验
type ArticlePageQueryInput struct {
	PageNo       int    `form:"page_no" binding:"omitempty,min=1"`          // min=1：页码必须大于等于 1
	PageSize     int    `form:"page_size" binding:"omitempty,min=1,max=50"` // max=50：限制单页最大返回 50 条，防止前端恶意传入 size=10000 导致数据库宕机
	CategoryName string `form:"category_name" binding:"omitempty"`          // 分类筛选是可选的
	OrderBy      string `form:"order_by" binding:"omitempty"`               // 排序
	State        *int8  `form:"state" binding:"omitempty,oneof=0 1"`        // 筛选状态，选填。用指针是为了能区分前端传的是 0 还是没传
}

// ArticleDetailInput 获取文章详情的请求入参
type ArticleDetailInput struct {
	ID uint64 `form:"id" uri:"id" binding:"required,gt=0"`
}

// ArticleDeleteInput 删除文章的请求入参
type ArticleDeleteInput struct {
	// 绑定并校验至少传一个 ID
	Ids []uint64 `json:"ids" binding:"required,gt=0" label:"文章ID列表"`
}
