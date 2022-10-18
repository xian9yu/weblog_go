package article

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

// 查询分类中的文章
func CategoryList(c *gin.Context) {
	categoryId := c.Query("category_id")               // 分类id
	pageNo := utils.AnyToInt(c.Query("page_no"))       // 页码
	limitSize := utils.AnyToInt(c.Query("limit_size")) // 每页显示数量
	orderBy := c.Query("order_by")                     // order by (id DESC)

	if limitSize <= 0 {
		limitSize = 10 // 每页显示数量默认为10
	}

	// string 转 []string
	var categoryIds []int
	_ = json.Unmarshal([]byte(categoryId), &categoryIds)

	// 通过分类 id 查询文章 id
	var links models.Links
	articleId, err := links.GetList(categoryIds)
	if err != nil {
		response.Error(c, "获取列表失败", err)
		c.Abort()
		return
	}
	if articleId == nil {
		response.Error(c, "查询分类为空", nil)
		c.Abort()
		return
	}

	// 取 article id
	var articleIds []int
	for _, id := range articleId {
		articleIds = append(articleIds, utils.AnyToInt(id.ArticleId))
	}

	// 通过 article id 查询文章列表
	var article models.Article
	total, list, err := article.GetCategoryList(articleIds, orderBy, utils.Pagination{
		PageNo:    pageNo,
		LimitSize: limitSize,
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
}

// 文章列表
func List(c *gin.Context) {
	pageNo := utils.AnyToInt(c.Query("page_no"))       // 页码
	limitSize := utils.AnyToInt(c.Query("limit_size")) // 每页显示数量
	orderBy := c.Query("order_by")                     // order by (id DESC)

	if limitSize == 0 {
		limitSize = 10 // 每页显示数量默认为10
	}

	// 通过 article id 查询文章列表
	var article models.Article
	total, list, err := article.GetList(orderBy, utils.Pagination{
		PageNo:    pageNo,
		LimitSize: limitSize,
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
}
