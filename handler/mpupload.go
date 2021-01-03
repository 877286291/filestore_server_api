package handler

import (
	rdb "github.com/aurora/Filestore-server/db/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"strconv"
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

}
