package dto

// AdminResetPasswordInput 管理员后台强制重置他人密码
type AdminResetPasswordInput struct {
	TargetUserId uint64 `json:"target_user_id" binding:"required" label:"目标用户ID"`
	NewPassword  string `json:"new_password" binding:"required,min=6,max=32" label:"新密码"`
}
