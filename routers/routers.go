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
		fileRouter.POST("/fastUpload", handler.FastUploadHandler)
		fileRouter.GET("/meta", handler.GetFileMetaHandler)
		fileRouter.GET("/list", handler.FileQueryHandler)
		fileRouter.GET("/download", handler.DownloadHandler)
		fileRouter.PUT("/update", handler.FileMetaUpdateHandler)
		fileRouter.DELETE("/delete", handler.FileDeleteHandler)
		mpUpload := fileRouter.Group("/mpupload")
		{
			mpUpload.POST("/init", handler.InitMultipartUploadHandler)
			mpUpload.POST("/uppart", handler.UploadPartHandler)
			mpUpload.POST("/complete", handler.CompleteUploadHandler)
			mpUpload.GET("/cancel", handler.CancelUploadHandler)
			mpUpload.GET("/uploadStatus", handler.MultipartUploadStatusHandler)
		}
	}
	userRouter := e.Group("/user")
	{
		userRouter.POST("/signUp", handler.SignUpHandler)
		userRouter.POST("/signIn", handler.SignInHandler)
		userRouter.GET("/info", middleware.JWT(), handler.UserInfoHandler)
	}
}
