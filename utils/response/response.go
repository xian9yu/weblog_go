package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一返回的 JSON 骨架（首字母大写，方便外部包和文档引用）
type Response struct {
	Code int `json:"code"` // 业务自定义错误码（可选）
	Msg  any `json:"msg"`  // 提示信息，用 any 可以接收 string 或结构化错误
	Data any `json:"data"` // 数据内容
}

// Result 核心底层方法：真正决定 HTTP 状态码的地方
func result(c *gin.Context, httpStatus int, businessCode int, msg any, data any) {
	c.JSON(httpStatus, Response{
		Code: businessCode,
		Msg:  msg,
		Data: data,
	})
}

// Ok 200 成功请求
func Ok(c *gin.Context, data any) {
	result(c, http.StatusOK, 200, "success", data)
}

// OkMsg 200 成功请求（带自定义成功文案，常用于操作类接口）
func OkMsg(c *gin.Context, msg string) {
	result(c, http.StatusOK, 200, msg, nil)
}

// Created 201 成功创建资源（符合标准 RESTful 规范）
func Created(c *gin.Context, data any) {
	result(c, http.StatusCreated, 201, "created", data)
}

// ====================  失败响应系列 ====================

// FailClient 400 客户端参数错误
// 支持格式化传参，例如: FailClient(c, "用户 %s 不存在", username)
func FailClient(c *gin.Context, format string, args ...any) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	result(c, http.StatusBadRequest, 400, msg, nil)
}

// FailUnauthorized 401 认证失败、未登录或 Token 过期
func FailUnauthorized(c *gin.Context, format string, args ...any) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	result(c, http.StatusUnauthorized, 401, msg, nil)
}

// FailForbidden 403 权限不足（比如普通用户想去删管理员的文章，新增常用状态）
func FailForbidden(c *gin.Context, format string, args ...any) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	result(c, http.StatusForbidden, 403, msg, nil)
}

// FailServer 500 服务器内部致命错误
func FailServer(c *gin.Context, format string, args ...any) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	result(c, http.StatusInternalServerError, 500, msg, nil)
}
