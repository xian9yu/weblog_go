package article

import (
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

func (ac *Controller) Delete(c *gin.Context) {
	var input dto.ArticleUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, "参数解析失败: %v", err)
		return
	}
	// 获取 token 详情
	currentUserId, _ := c.Get("user_id")
	group, _ := c.Get("group")

	// 判断 如果不是超级管理员，就必须进行严苛的“作者本人”所有权校验
	if group != "admin" {
		// 获取 article 的数据
		article, err := ac.articleRepo.GetArticleDetailsById(input.ID)
		if err != nil || article == nil {
			response.FailClient(c, "删除失败，目标文章不存在")
			return
		}

		// 检查这篇文章的作者 ID，和当前登录的用户 UserId 是否一致
		if article.UserId != currentUserId {
			response.FailForbidden(c, "对不起，您没有权限删除他人的文章")
			return
		}
	}

	rowsAffected := ac.articleRepo.Delete(input.ID)
	if rowsAffected != 1 {
		response.FailClient(c, "服务器繁忙，文章删除失败，请稍后再试")
		return
	}

	response.OkMsg(c, "文章删除成功")
}
