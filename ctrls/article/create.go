package article

import (
	"strings"
	"time"
	"weblog/dto"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"

	"github.com/gin-gonic/gin"
)

func (ac *Controller) Create(c *gin.Context) {
	var input dto.ArticleCreateInput

	// c.ShouldBindJSON 会自动根据 binding 标签进行校验
	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果前端传参不合法（比如 title 没传，或者 status 传了 3）
		// 这里会直接拦截，并返回具体的错误原因
		response.FailClient(c, "参数校验失败", err.Error())
		return
	}

	// 判断 title 是否重复
	// 新建的文章在数据库里还没有 ID，所以排除 ID 传 0 即可
	// 先修剪首尾空格
	input.Title = strings.TrimSpace(input.Title)
	exists, err := ac.articleRepo.ArticleTitleExists(input.Title, 0)
	if err != nil {
		response.FailServer(c, "系统繁忙，查重失败")
		return
	}
	if exists {
		response.FailClient(c, "文章标题「%s」已存在，请更换", input.Title)
		return
	}

	// 获取 token 详情
	userId, _ := c.Get("user_id")

	nowTime := uint64(time.Now().Unix())
	article := models.Article{
		Title:       input.Title,
		Content:     input.Content,
		UserId:      uint64(utils.AnyToInt(userId)),
		Status:      input.Status, // 0 代表草稿，1 代表已发布。传 2 或 -1 直接报错拦截
		CreatedTime: nowTime,
		UpdatedTime: nowTime,
	}

	articleId, err := ac.articleRepo.Add(article)
	if err != nil {
		response.FailClient(c, "保存文章失败", err.Error())
	}
	response.OkMsg(c, "文章保存成功，article_id: "+utils.AnyToString(articleId))
}
