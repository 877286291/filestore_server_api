package routers

import (
	"github.com/aurora/Filestore-server/handler"
	"github.com/aurora/Filestore-server/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouters(e *gin.Engine) {
	fileRouter := e.Group("/file")
	fileRouter.Use(middleware.JWT())
	{
		fileRouter.POST("/upload", handler.UploadHandler)
		fileRouter.GET("/meta", handler.GetFileMetaHandler)
		fileRouter.GET("/list", handler.FileQueryHandler)
		fileRouter.GET("/download", handler.DownloadHandler)
		fileRouter.PUT("/update", handler.FileMetaUpdateHandler)
		fileRouter.DELETE("/delete", handler.FileDeleteHandler)
	}
	userRouter := e.Group("/user")
	{
		userRouter.POST("/signUp", handler.SignUpHandler)
		userRouter.POST("/signIn", handler.SignInHandler)
	}
}
