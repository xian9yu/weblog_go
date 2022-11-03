package user

import (
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

// 查询 用户列表
func List(c *gin.Context) {
	group := c.Query("group")                         // 分组 id
	pageNo := utils.AnyToInt(c.Query("page_no"))      // 页码
	pageSize := utils.AnyToInt(c.Query("limit_size")) // 每页显示数量
	orderBy := c.Query("order_by")                    // order by (id DESC)

	// 用户未输入页码, 则默认为 1
	if pageNo == 0 {
		pageNo = 1
	}

	// 获取 token 详情
	tokenGroup, _ := c.Get("group")
	// 只有 管理员 才能操作
	if tokenGroup == "admin" {
		// 查询用户列表
		var user models.User
		total, list, err := user.GetList(group, orderBy, utils.Pagination{
			PageNo:   pageNo,
			PageSize: utils.Pagination{}.DefaultPageSize(pageSize), // 判断每页数量是否为空, 为空则给默认值
		})
		if err != nil {
			response.Error(c, "获取列表失败", err)
			c.Abort()
			return
		}

		response.Success(c, gin.H{
			"total":   total,
			"page_no": pageNo,
			"list":    list,
		})
		c.Abort()
		return
	}

	response.PageNotFound(c)
}
