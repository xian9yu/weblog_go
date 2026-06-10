package article

import (
	"errors"
	"strings"
	"time"
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (ac *Controller) Update(c *gin.Context) {
	var input dto.ArticleUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "文章更新参数格式不正确")
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	exists, err := ac.repos.Article.ArticleTitleExists(input.Title, input.ID)
	if err != nil {
		response.FailServer(c, "系统繁忙，标题校验失败")
		return
	}
	if exists {
		response.FailClient(c, "文章标题「%s」已存在，请更换标题", input.Title)
		return
	}

	// 获取 token 详情
	currentUserId, _ := c.Get("user_id")
	group, _ := c.Get("group")

	// 先修剪首尾空格
	if group != "admin" {
		// 获取 authorId
		authorId, err := ac.repos.Article.GetArticleUserIdById(input.ID)
		if err != nil {
			// 如果是 GORM 没查到记录，通常会报 ErrRecordNotFound
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.FailClient(c, "操作失败，目标文章不存在")
				return
			}
			response.FailServer(c, "系统繁忙，权限校验失败")
			return
		}

		if authorId != currentUserId {
			response.FailForbidden(c, "对不起，您没有权限操作他人的文章")
			return
		}
	}

	updateMap := map[string]any{
		"title":      input.Title,
		"content":    input.Content,
		"status":     input.Status,
		"updated_at": time.Now(),
	}
	rows, err := ac.repos.Article.Update(input.ID, updateMap)
	if err != nil {
		response.FailServer(c, "系统繁忙，文章更新失败")
		return
	}

	// 判断是否真正影响了行数
	if rows == 0 {
		response.FailClient(c, "文章未做任何修改或文章不存在")
		return
	}

	response.OkMsg(c, "文章更新成功")
}
