package dto

// AdminResetPasswordInput 管理员后台强制重置他人密码
type AdminResetPasswordInput struct {
	TargetUserId uint64 `json:"target_user_id" binding:"required" label:"目标用户ID"`
	NewPassword  string `json:"new_password" binding:"required,min=6,max=32" label:"新密码"`
}

// UserListInput 用户列表查询入参 DTO
type UserListInput struct {
	Group    string `form:"group"`
	PageNo   int    `form:"page_no" binding:"omitempty,min=1"`          // min=1：页码必须大于等于 1
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=50"` // max=50：限制单页最大返回 50 条，防止前端恶意传入 size=10000 导致数据库宕机
	OrderBy  string `form:"order_by"`
}
