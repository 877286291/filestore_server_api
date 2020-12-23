package handler

import (
	"github.com/aurora/Filestore-server/meta"
	"github.com/aurora/Filestore-server/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
		return
	}
	fileMeta := meta.FileMeta{
		FileSha1: utils.FileSha1(newFile),
		FileName: file.Filename,
		FileSize: utils.FileSizeConversion(int(file.Size)),
		Location: file.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	meta.UpdateFileMeta(fileMeta)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "文件上传成功", "size": utils.FileSizeConversion(int(file.Size))})
}
func GetFileMetaHandler(c *gin.Context) {
	fileHash := c.DefaultQuery("filehash", "")
	if fileHash == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "请检查请求参数"})
		return
	}
	fileMeta := meta.GetFileMeta(fileHash)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "文件信息查询成功", "data": fileMeta})
}
func FileQueryHandler(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	ErrHandler(err)
	metas := meta.GetListFileMetas(limit)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "查询成功", "data": metas})
}
func DownloadHandler(c *gin.Context) {
	fileSha1 := c.DefaultQuery("filehash", "")
	fileMeta := meta.GetFileMeta(fileSha1)
	if fileMeta.FileSha1 == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileMeta.FileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(fileMeta.Location)
}
func FileMetaUpdateHandler(c *gin.Context) {
	opType := c.DefaultQuery("op", "0")
	fileSha1 := c.DefaultQuery("filehash", "")
	newFileName := c.DefaultQuery("filename", "")
	if opType != "0" {
		c.IndentedJSON(http.StatusForbidden, gin.H{"msg": ""})
		return
	}
	curFileMeta := meta.GetFileMeta(fileSha1)
	if curFileMeta.UploadAt == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "文件信息更新成功", "data": curFileMeta})
}

func FileDeleteHandler(c *gin.Context) {
	fileSha1 := c.DefaultQuery("filehash", "")
	fileMeta := meta.GetFileMeta(fileSha1)
	if fileMeta.FileSha1 == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := os.Remove(fileMeta.Location)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "文件删除失败!"})
		return
	}
	meta.RemoveFileMeta(fileSha1)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "文件删除成功!"})
}
