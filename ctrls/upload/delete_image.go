package upload

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

func Delete(c *gin.Context) {
	id := c.PostForm("id")
	fileId := uint64(utils.AnyToInt(id))

	// 事务开始
	var db = models.MySQL()
	err := db.Transaction(func(tx *gorm.DB) error {
		// 获取文件详情
		var upload models.Upload
		upload.Id = fileId
		details, err := upload.TxGetDetailsById(tx)
		if err != nil {
			return err
		}

		if details.Id < 1 {
			return errors.New("资源不存在")
		}

		// 获取 token 详情
		userId, _ := c.Get("user_id")
		group, _ := c.Get("group")
		// 只有 (本人/管理员) 才能操作
		if details.UserId == uint64(utils.AnyToInt(userId)) || group == "admin" {
			// 删除数据库
			rowsAffected := upload.TxDelete(tx)
			if rowsAffected != 1 {
				return err
			}

			// 从 服务器 删除文件
			if os.Remove(details.FullUrl) != nil {
				return err
			}
			return nil
		}
		return err
	})
	if err != nil {
		response.Error(c, "删除失败, 请稍后再试", err)
		c.Abort()
		return
	}

	response.Success(c, gin.H{
		"message": "删除成功",
	})
}
