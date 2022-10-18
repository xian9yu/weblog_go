package upload

import (
	"github.com/gin-gonic/gin"
	"path"
	"strings"
	"time"
	"weblog/models"
	"weblog/utils"
	"weblog/utils/response"
)

// 单文件上传 (上传到本地文件夹)
func Image(c *gin.Context) {
	f, _ := c.FormFile("image")

	// 获取 path 路径
	v := utils.InitConfig()
	cPath := v.GetString("upload.image.path")

	// 分割文件名 || 文件格式
	dst := cPath + f.Filename

	// 上传文件夹不存在则创建
	err := utils.CreateDir(cPath)
	if err != nil {
		response.Error(c, "文件夹创建失败", err)
		c.Abort()
		return
	}

	if models.CountUploadByAny("name", f.Filename) > 0 {
		response.Error(c, "上传失败, 已存在同名文件, 请修改后重试", nil)
		c.Abort()
		return
	}
	// 限制上传图片后缀
	fileExt := strings.ToLower(path.Ext(f.Filename))
	switch fileExt {
	case ".jpg":
	case ".jpeg":
	case ".png":
	case ".gif":
	default:
		response.Error(c, "仅支持 jpg,jpeg,png,gif 格式图片", nil)
		c.Abort()
		return
	}

	// 插入数据库
	upload := models.Upload{
		Name:        f.Filename,
		FullName:    f.Filename + fileExt,
		Path:        cPath,
		FullUrl:     dst,
		Size:        uint64(f.Size),
		Suffix:      fileExt,
		Category:    "image",
		CreatedTime: uint64(time.Now().Unix()),
	}
	id, rowsAffected, er := upload.Add()
	if er != nil || rowsAffected != 1 || id < 1 {
		response.Error(c, "上传文件失败, 请重试", er)
		c.Abort()
		return
	}

	// 上传文件至指定的完整文件路径 ( 路径必须已存在 )
	err = c.SaveUploadedFile(f, dst)
	if err == nil {
		response.Success(c, gin.H{
			"message":  "上传成功",
			"path":     dst,
			"image_id": id,
		})
		c.Abort()
		return
	}

	response.Error(c, "上传文件失败, 请稍后再试", err)
}
