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
