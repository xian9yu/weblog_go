package dto

// UserRegisterInput 用户注册输入参数
type UserRegisterInput struct {
	Password string `json:"password" binding:"required,min=6,max=32" label:"密码"`
	Name     string `json:"name" binding:"required,max=50" label:"昵称"`
	Email    string `json:"email" binding:"required,email,max=100" label:"邮箱"`
}

// UserLoginInput 用户登录输入参数
type UserLoginInput struct {
	Email    string `json:"email" binding:"omitempty,email,max=100" label:"邮箱"`
	Password string `json:"password" binding:"required" label:"密码"`
}

// UserUpdateProfileInput 作者在后台修改个人资料（不含密码）
type UserUpdateProfileInput struct {
	Name  string `json:"name" binding:"required,max=50" label:"昵称"`
	Email string `json:"email" binding:"omitempty,email,max=100" label:"邮箱"`
	// 如果以后有头像、简介等字段，可以在这里直接扩充
}

// UserUpdatePasswordInput 修改密码输入参数
type UserUpdatePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required" label:"旧密码"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32" label:"新密码"`
}

// UserListInput 用户列表查询入参 DTO
type UserListInput struct {
	Group    string `form:"group"`
	PageNo   int    `form:"page_no" binding:"omitempty,min=1"`          // min=1：页码必须大于等于 1
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=50"` // max=50：限制单页最大返回 50 条，防止前端恶意传入 size=10000 导致数据库宕机
	OrderBy  string `form:"order_by"`
}

// UserDetailInput 获取用户详情的请求入参
type UserDetailInput struct {
	ID    uint64 `uri:"id" binding:"required_without=Email,omitempty,gt=0" label:"用户ID"`
	Email string `uri:"email" binding:"required_without=ID,omitempty,email,max=100" label:"用户邮箱"`
}

// UserUnregisterInput 用户自主注销入参
type UserUnregisterInput struct {
	Password string `json:"password" binding:"required" label:"登录密码"`
}
