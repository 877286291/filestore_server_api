package routers

import (
	"github.com/aurora/Filestore-server/handler"
	"github.com/gin-gonic/gin"
)

func InitRouters(e *gin.Engine) {
	e.POST("/file/upload", handler.UploadHandler)
	e.GET("/file/meta", handler.GetFileMetaHandler)
	e.GET("/file/list", handler.FileQueryHandler)
	e.GET("/file/download", handler.DownloadHandler)
	e.DELETE("/file/delete", handler.DeleteHandler)
}
