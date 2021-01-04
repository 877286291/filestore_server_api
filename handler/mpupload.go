package handler

import (
	"github.com/aurora/Filestore-server/db"
	rdb "github.com/aurora/Filestore-server/db/redis"
	"github.com/aurora/Filestore-server/meta"
	"github.com/aurora/Filestore-server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MultipartUploadInfo struct {
	FileHash   string
	FileName   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func InitMultipartUploadHandler(c *gin.Context) {
	fileHash := c.PostForm("filehash")
	filename := c.PostForm("filename")
	filesize, _ := strconv.Atoi(c.PostForm("filesize"))
	rConn := rdb.RedisConn()
	upInfo := MultipartUploadInfo{
		FileHash:   fileHash,
		FileName:   filename,
		FileSize:   filesize,
		UploadID:   uuid.New().String(),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize / (5 * 1024 * 1024)))),
	}
	rConn.HMSet("MP_"+upInfo.UploadID, map[string]interface{}{
		"chunkCount": upInfo.ChunkCount,
		"fileHash":   upInfo.FileHash,
		"FileName":   upInfo.FileName,
		"FileSize":   upInfo.FileSize,
	})
	c.JSON(http.StatusOK, gin.H{"msg": "分块上传文件信息初始化成功"})
}
func UploadPartHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")
	chunkIndex := c.PostForm("index")
	blockFile, err := c.FormFile("blockfile")
	rConn := rdb.RedisConn()
	filePath := "data/" + uploadID + "/" + chunkIndex
	_ = os.MkdirAll(path.Dir(filePath), 0777)
	fd, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建文件请求被拒绝，权限不足"})
		return
	}
	defer fd.Close()
	content, err := blockFile.Open()
	_, _ = io.Copy(fd, content)
	rConn.HSet("MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	c.JSON(http.StatusOK, gin.H{"msg": "chkidx_" + chunkIndex + "上传完成"})
}
func CompleteUploadHandler(c *gin.Context) {
	claims, _ := utils.ParseToken(c.GetHeader("token"))
	uploadID := c.PostForm("uploadid")
	fileHash := c.PostForm("filehash")
	filename := c.PostForm("filename")
	fileSize, _ := strconv.Atoi(c.PostForm("filesize"))
	rConn := rdb.RedisConn()
	data := rConn.HGetAll("MP_" + uploadID).Val()
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
	fileMeta := meta.FileMeta{
		FileSha1: fileHash,
		FileName: filename,
		FileSize: int64(fileSize),
		Location: "",
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	if meta.UpdateFileMetaDB(fileMeta) && db.UserFileUploaded(claims.Username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize) {
		c.JSON(http.StatusOK, gin.H{"msg": "文件上传成功", "size": utils.FileSizeConversion(fileSize)})
		return
	}
}
