package handler

import (
	"github.com/aurora/Filestore-server/db"
	"github.com/aurora/Filestore-server/meta"
	"github.com/aurora/Filestore-server/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(c *gin.Context) {
	claims, _ := utils.ParseToken(c.GetHeader("token"))
	file, err := c.FormFile("file")
	ErrHandler(err)
	newFile, err := os.Create(file.Filename)
	ErrHandler(err)
	defer newFile.Close()
	content, err := file.Open()
	ErrHandler(err)
	defer content.Close()
	written, err := io.Copy(newFile, content)
	_, _ = newFile.Seek(0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件上传失败", "size": utils.FileSizeConversion(int(written))})
		return
	}
	fileMeta := meta.FileMeta{
		FileSha1: utils.FileSha1(newFile),
		FileName: file.Filename,
		FileSize: file.Size,
		Location: file.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	if meta.UpdateFileMetaDB(fileMeta) && db.UserFileUploaded(claims.Username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize) {
		c.JSON(http.StatusOK, gin.H{"msg": "文件上传成功", "size": utils.FileSizeConversion(int(file.Size))})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件上传失败"})
}
func GetFileMetaHandler(c *gin.Context) {
	fileHash := c.DefaultQuery("filehash", "")
	if fileHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请检查请求参数"})
		return
	}
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件信息获取失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "文件信息查询成功", "data": fileMeta})
}
func FileQueryHandler(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	ErrHandler(err)
	claims, _ := utils.ParseToken(c.GetHeader("token"))
	metas, err := db.GetUserFileMetas(claims.Username, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件信息获取失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "文件信息查询成功", "data": metas})
}
func DownloadHandler(c *gin.Context) {
	fileSha1 := c.DefaultQuery("filehash", "")
	if fileSha1 == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
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
	if opType != "0" || fileSha1 == "" || newFileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请检查请求参数"})
		return
	}
	curFileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件信息更新失败"})
		return
	}
	curFileMeta.FileName = newFileName
	_ = meta.UpdateFileMetaDB(curFileMeta)
	c.JSON(http.StatusOK, gin.H{"msg": "文件信息更新成功", "data": curFileMeta})
}

func FileDeleteHandler(c *gin.Context) {
	fileSha1 := c.DefaultQuery("filehash", "")
	fileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件删除失败!"})
		return
	}
	_ = os.Remove(fileMeta.Location)
	if meta.RemoveFileMetaDB(fileSha1) {
		c.JSON(http.StatusOK, gin.H{"msg": "文件删除成功!"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件删除失败!"})
}
func FastUploadHandler(c *gin.Context) {
	claims, _ := utils.ParseToken(c.GetHeader("token"))
	fileHash := c.PostForm("filehash")
	filename := c.PostForm("filename")
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil || fileMeta.FileSha1 == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "秒传失败"})
		return
	}
	if db.UserFileUploaded(claims.Username, fileMeta.FileSha1, filename, fileMeta.FileSize) {
		fileMeta.FileName = filename
		c.JSON(http.StatusOK, gin.H{"msg": "秒传成功", "data": fileMeta})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"msg": "秒传失败，请稍后重试"})
}
