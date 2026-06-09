package dto

// ConfigCreateInput 新增配置项的入参
// 适用于路由：POST /api/v1/admin/config
type ConfigCreateInput struct {
	Key string `json:"key" binding:"required,min=2,max=60" label:"配置键名"`
	// Value 允许为空字符串（比如有些开关默认是关闭的），但必须限制最大长度防止撑爆数据库
	Value string `json:"value" binding:"max=65535" label:"配置项参数值"`
	// Remark 强烈建议必填，必须让管理员写清楚这个配置是干嘛用的，方便后续维护
	Remark string `json:"remark" binding:"required,max=250" label:"备注介绍"`
}

// ConfigGetByKeyInput 获取单个配置项的入参
// 适用于路由：GET /api/v1/admin/config/:key (路径参数版)
type ConfigGetByKeyInput struct {
	// 绑定路径中的 :key，强校验不能为空，且限制长度防止恶意注入
	Key string `uri:"key" binding:"required,max=60" label:"配置键名"`
}

// ConfigGetByIdInput 通过ID查询配置的入参
// 路由：GET /api/v1/admin/config/:id
type ConfigGetByIdInput struct {
	ID uint64 `uri:"id" binding:"required,gt=0" label:"配置ID"`
}

// ConfigUpdateInput 修改配置项的入参
// 适用于路由：PUT /api/v1/admin/config
type ConfigUpdateInput struct {
	Key    string `json:"key" binding:"required,max=60" label:"配置键名"`
	Value  string `json:"value" binding:"max=65535" label:"配置项参数值"`
	Remark string `json:"remark" binding:"max=250" label:"备注介绍"`
}

// ConfigListInput 获取配置项列表的入参（查询参数）
// 适用于路由：GET /api/v1/admin/config
type ConfigListInput struct {
	PageNo   int `form:"page_no" binding:"required,gt=0" label:"当前页码"`
	PageSize int `form:"page_size" binding:"required,gt=0,max=100" label:"每页条数"`
	// Keyword 模糊搜索关键字（可以支持按 key 或 remark 搜索），非必填
	Keyword string `form:"keyword" binding:"max=50" label:"搜索关键字"`
}

// ConfigDeleteInput 删除系统配置项入参
type ConfigDeleteInput struct {
	// Ids 支持批量删除和单个删除通用
	// binding:"required" 确保前端必须传这个字段
	// binding:"gt=0" 确保数组内至少有一个 ID，防止传空数组 [] 导致无效请求
	Ids []uint64 `json:"ids" binding:"required,gt=0"`
}
