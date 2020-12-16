package handler

import (
	"github.com/aurora/Filestore-server/meta"
	"github.com/aurora/Filestore-server/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	ErrHandler(err)
	log.Println(file.Filename + ":" + utils.FileSizeConversion(int(file.Size)))
	newFile, err := os.Create(file.Filename)
	ErrHandler(err)
	defer newFile.Close()
	content, err := file.Open()
	ErrHandler(err)
	defer content.Close()
	written, err := io.Copy(newFile, content)
	_, _ = newFile.Seek(0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "文件上传失败", "size": utils.FileSizeConversion(int(written))})
	} else {
		fileMeta := meta.FileMeta{
			FileSha1: utils.FileSha1(newFile),
			FileName: file.Filename,
			FileSize: file.Size,
			Location: file.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		meta.UpdateFileMeta(fileMeta)
		c.IndentedJSON(http.StatusOK, gin.H{"msg": "文件上传成功", "size": utils.FileSizeConversion(int(file.Size))})
	}
}
func GetFileMetaHandler(c *gin.Context) {
	fileHash := c.DefaultQuery("filehash", "")
	fileMeta := meta.GetFileMeta(fileHash)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "文件信息查询成功", "data": fileMeta})
}
func FileQueryHandler(c *gin.Context) {
	meta.GetListFileMetas(1)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "查询成功"})
}
