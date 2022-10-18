package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
func Error(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    -1,
		"message": message,
		"error":   err,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code":    -1,
		"message": message,
	})
}

func PageNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"code":    -1,
		"message": "404 PAGE NOT FOUND",
	})
}
