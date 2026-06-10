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
			c.Abort()
			return
		}
	}
	// 如果前端没传 page_no，或者传了 0，默认第 1 页
	if input.PageNo <= 0 {
		input.PageNo = 1
	}

	// 如果前端没传 page_size，或者传了 0，默认每页 10 条
	if input.PageSize <= 0 {
		input.PageSize = 10
	}

	total, list, err := ac.repos.Article.GetList(input, 0)
	if err != nil {
		response.FailServer(c, "系统繁忙，获取文章列表失败")
		c.Abort()
		return
	}

	response.Ok(c, gin.H{
		"total": total,
		"list":  list,
		"page":  input.PageNo,
	})
}

// AdminList 路由大概长这样：GET /api/v1/admin/categories/:name/articles （需要 JWT 中间件拦截）
func (ac *Controller) AdminList(c *gin.Context) {
	var input dto.ArticlePageQueryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if err := c.ShouldBindJSON(&input); err != nil {
			response.FailClient(c, "参数解析失败: %v", err)
			c.Abort()
			return
		}
	}

	// 如果前端没传 page_no，或者传了 0，默认第 1 页
	if input.PageNo <= 0 {
		input.PageNo = 1
	}

	// 如果前端没传 page_size，或者传了 0，默认每页 10 条
	if input.PageSize <= 0 {
		input.PageSize = 10
	}

	currentUserId, exists := c.Get("userId")
	if !exists {
		response.FailUnauthorized(c, "登录已失效，请重新登录")
		c.Abort()
		return
	}

	total, list, err := ac.repos.Article.GetList(input, currentUserId.(uint64))
	if err != nil {
		response.FailServer(c, "系统繁忙，获取文章列表失败")
		c.Abort()
		return
	}

	response.Ok(c, gin.H{
		"total": total,
		"list":  list,
		"page":  input.PageNo,
	})
}
