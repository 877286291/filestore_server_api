package handler

import (
	"fmt"
	"github.com/aurora/Filestore-server/db"
	rdb "github.com/aurora/Filestore-server/db/redis"
	"github.com/aurora/Filestore-server/meta"
	"github.com/aurora/Filestore-server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MultipartUploadInfo struct {
	FileHash    string
	FileName    string
	FileSize    int
	UploadID    string
	ChunkSize   int
	ChunkCount  int
	ChunkExists []int
}

const (
	FilePath          = "data"
	ChunkKeyPrefix    = "MP_"
	HashUpIDKeyPrefix = "HASH_UPID_"
)

func InitMultipartUploadHandler(c *gin.Context) {
	fileHash := c.PostForm("filehash")
	filename := c.PostForm("filename")
	filesize, _ := strconv.Atoi(c.PostForm("filesize"))
	rConn := rdb.RedisConn()
	upInfo := MultipartUploadInfo{
		FileHash:    fileHash,
		FileName:    filename,
		FileSize:    filesize,
		UploadID:    uuid.New().String(),
		ChunkSize:   5 * 1024 * 1024,
		ChunkCount:  int(math.Ceil(float64(filesize / (5 * 1024 * 1024)))),
		ChunkExists: make([]int, 0),
	}
	rConn.HMSet(ChunkKeyPrefix+upInfo.UploadID, map[string]interface{}{
		"chunkCount": upInfo.ChunkCount,
		"fileHash":   upInfo.FileHash,
		"FileName":   upInfo.FileName,
		"FileSize":   upInfo.FileSize,
	})
	pwd, _ := os.Getwd()
	_ = os.MkdirAll(filepath.Join(pwd, FilePath, upInfo.UploadID), 0777)
	c.JSON(http.StatusOK, gin.H{"msg": "分块上传文件信息初始化成功", "data": upInfo})
}
func UploadPartHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")
	chunkIndex := c.PostForm("index")
	blockFile, err := c.FormFile("blockfile")
	rConn := rdb.RedisConn()
	fd, err := os.Create(filepath.Join(FilePath, uploadID, chunkIndex))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建文件请求被拒绝，权限不足"})
		return
	}
	defer fd.Close()
	content, err := blockFile.Open()
	_, _ = io.Copy(fd, content)
	rConn.HSet(ChunkKeyPrefix+uploadID, "chkidx_"+chunkIndex, 1)
	c.JSON(http.StatusOK, gin.H{"msg": "chkidx_" + chunkIndex + "上传完成"})
}
func CompleteUploadHandler(c *gin.Context) {
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	uploadID := c.PostForm("uploadid")
	fileHash := c.PostForm("filehash")
	filename := c.PostForm("filename")
	fileSize, _ := strconv.Atoi(c.PostForm("filesize"))
	rConn := rdb.RedisConn()
	data := rConn.HGetAll(ChunkKeyPrefix + uploadID).Val()
	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{"msg": "合并失败"})
	}
	totalCount, chunkCount := 0, 0
	for k, v := range data {
		if k == "chunkCount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount += 1
		}
	}
	if totalCount != chunkCount || totalCount == 0 || chunkCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部服务器错误"})
		return
	}
	for i := 1; i <= totalCount; i++ {
		src, _ := ioutil.ReadFile(filepath.Join(FilePath, uploadID, strconv.Itoa(i)))
		dst, _ := os.OpenFile(filepath.Join(FilePath, uploadID, filename), os.O_CREATE|os.O_WRONLY, os.ModePerm)
		_, err := dst.Write(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件合并失败"})
			return
		}
		_ = os.Remove(filepath.Join(FilePath, uploadID, strconv.Itoa(i)))
	}
	fileMeta := meta.FileMeta{
		FileSha1: fileHash,
		FileName: filename,
		FileSize: int64(fileSize),
		Location: "",
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	if meta.UpdateFileMetaDB(fileMeta) && db.UserFileUploaded(claims.Username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize) {
		c.JSON(http.StatusOK, gin.H{"msg": "文件上传成功", "size": utils.FileSizeConversion(fileSize)})
		rConn.Del(ChunkKeyPrefix + uploadID)
		return
	}
}
func CancelUploadHandler(c *gin.Context) {
	uploadId := c.Query("uploadid")
	err := os.RemoveAll(filepath.Join(FilePath, uploadId))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"msg": "上传任务取消失败"})
		return
	}
	rConn := rdb.RedisConn()
	rConn.Del(ChunkKeyPrefix + uploadId)
	c.JSON(http.StatusOK, gin.H{"msg": "上传任务取消成功"})
}
func MultipartUploadStatusHandler(c *gin.Context) {
	uploadId := c.Query("uploadid")
	rConn := rdb.RedisConn()
	data := rConn.HGetAll(ChunkKeyPrefix + uploadId).Val()
	chunkCount := 0
	for k, v := range data {
		fmt.Println(k, v)
		if strings.HasPrefix(k, "chkidx_") {
			chunkCount++
		}
	}
	totalCount, _ := strconv.Atoi(data["chunkCount"])
	data["percentage"] = fmt.Sprintf("%.2f", (float64(chunkCount)/float64(totalCount))*100) + "%"
	c.JSON(http.StatusOK, gin.H{"msg": "获取上传状态信息成功", "data": data})
}
