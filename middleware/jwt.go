package middleware

import (
	"github.com/aurora/Filestore-server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = http.StatusOK
		token := c.GetHeader("token")
		if token == "" {
			code = http.StatusBadRequest
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil || time.Now().Unix() > claims.ExpiresAt {
				code = http.StatusUnauthorized
			}
		}
		if code != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "请求无效，请重新登录",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
