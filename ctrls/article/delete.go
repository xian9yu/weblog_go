package article

import (
	"fmt"
	"weblog/dto"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

// BatchDelete 批量删除文章
// 路由：DELETE /api/v1/admin/article
func (ac *Controller) BatchDelete(c *gin.Context) {
	var input dto.ArticleDeleteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailClient(c, err.Error())
		c.Abort()
		return
	}

	currentUserId := c.GetUint64("user_id")

	// 调用 Repo 执行安全删除
	rowsAffected, err := ac.repos.Article.BatchDelete(input.Ids, currentUserId)
	if err != nil {
		response.FailServer(c, "批量删除文章失败")
		c.Abort()
		return
	}

	// 极致的反馈体验
	if rowsAffected == int64(len(input.Ids)) {
		response.OkMsg(c, "选中的文章已成功删除")
	} else {
		// 可能是因为其中有几篇不是他写的，被 Repo 里的 user_id = currentUserId 拦截了
		response.OkMsg(c, fmt.Sprintf("操作完成，成功删除 %d 篇文章，其余文章因权限不足或不存在未处理", rowsAffected))
	}
}

// 单文章删除
//func (ac *Controller) Delete(c *gin.Context) {
//	var input dto.ArticleDeleteInput
//	if err := c.ShouldBindJSON(&input); err != nil {
//		response.FailClient(c, "参数解析失败: %v", err)
//		return
//	}
//	// 获取 token 详情
//	currentUserId, _ := c.Get("user_id")
//	group, _ := c.Get("group")
//
//	// 判断 如果不是超级管理员，就必须进行严苛的“作者本人”所有权校验
//	if group != "admin" {
//		// 获取 article 的数据
//		article, err := ac.articleRepo.GetArticleInfosById(input.ID)
//		if err != nil || article == nil {
//			response.FailClient(c, "删除失败，目标文章不存在")
//			return
//		}
//
//		// 检查这篇文章的作者 ID，和当前登录的用户 UserId 是否一致
//		if article.UserId != currentUserId {
//			response.FailForbidden(c, "对不起，您没有权限删除他人的文章")
//			return
//		}
//	}
//
//	rowsAffected := ac.articleRepo.Delete(input.ID)
//	if rowsAffected != 1 {
//		response.FailClient(c, "服务器繁忙，文章删除失败，请稍后再试")
//		return
//	}
//
//	response.OkMsg(c, "文章删除成功")
//}
