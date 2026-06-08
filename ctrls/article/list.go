package article

import (
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// GuestList 公开列表
func (ac *Controller) GuestList(c *gin.Context) {
	// 路由大概长这样：GET /api/v1/categories/:name/articles

	var input dto.ArticlePageQueryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if err := c.ShouldBindJSON(&input); err != nil {
			response.FailClient(c, "参数解析失败: %v", err)
			return
		}
	}
	// 2. 规范化分页默认值，防止前端没传导致计算出负数
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}

	total, list, err := ac.articleRepo.GetList(input, 0)
	if err != nil {
		response.FailServer(c, "系统繁忙，获取文章列表失败")
		return
	}

	response.Ok(c, gin.H{
		"total": total,
		"list":  list,
		"page":  input.Page,
	})
}

// 路由大概长这样：GET /api/v1/admin/categories/:name/articles （需要 JWT 中间件拦截）
func (ac *Controller) AdminList(c *gin.Context) {
	var input dto.ArticlePageQueryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if err := c.ShouldBindJSON(&input); err != nil {
			response.FailClient(c, "参数解析失败: %v", err)
			return
		}
	}

	// 2. 规范化分页默认值，防止前端没传导致计算出负数
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}

	currentUserId, exists := c.Get("userId")
	if !exists {
		response.FailUnauthorized(c, "登录已失效，请重新登录")
		return
	}

	total, list, err := ac.articleRepo.GetList(input, currentUserId.(uint64))
	if err != nil {
		response.FailServer(c, "系统繁忙，获取文章列表失败")
		return
	}

	response.Ok(c, gin.H{
		"total": total,
		"list":  list,
		"page":  input.Page,
	})
}
