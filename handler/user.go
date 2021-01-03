package handler

import (
	"github.com/aurora/Filestore-server/db"
	"github.com/aurora/Filestore-server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUpHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if len(username) < 3 || len(password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请检查用户参数"})
		return
	}
	encPasswd := utils.Sha1([]byte(password))
	if suc := db.UserSignUp(username, encPasswd); suc {
		c.JSON(http.StatusOK, gin.H{"msg": "用户注册成功"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"msg": "用户注册失败"})
}
func SignInHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	encPassword := utils.Sha1([]byte(password))
	if ok, token := db.UserSignIn(username, encPassword); ok {
		c.JSON(http.StatusOK, gin.H{"msg": "登录成功", "token": token})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"msg": "登录失败！"})
}
func UserInfoHandler(c *gin.Context) {

	username := c.Query("username")
	token := c.GetHeader("token")
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "认证失败"})
		return
	}
	if username == parseToken.Username {
		userInfo, err := db.GetUserInfo(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部服务器错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "用户信息查询成功", "data": userInfo})
		return
	}
	c.JSON(http.StatusForbidden, gin.H{"msg": "非法请求"})
}
