package meta

import (
	"github.com/aurora/Filestore-server/db"
	"time"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}
func UpdateFileMetaDB(meta FileMeta) bool {
	return db.OnnFileUploadFinished(meta.FileSha1, meta.FileName, meta.FileSize, meta.Location)
}

func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	return FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
		UploadAt: tfile.CreateAt.Time.Format(time.RFC3339),
	}, nil
}
func SortByUploadTime(fileMetas []FileMeta) []FileMeta {
	for index, v := range fileMetas {
		for i := index + 1; i < len(fileMetas); i++ {
			if v.UploadAt < fileMetas[i].UploadAt {
				fileMetas[index], fileMetas[i] = fileMetas[i], fileMetas[index]
			}
		}
	}
	return fileMetas
}
func GetListFileMetas(count int) []FileMeta {
	metaArray := make([]FileMeta, 0)
	for _, v := range fileMetas {
		metaArray = append(metaArray, v)
	}
	SortByUploadTime(metaArray)
	if count > len(metaArray) {
		return metaArray
	}
	return metaArray[0:count]
}
func RemoveFileMetaDB(fileSha1 string) bool {
	return db.DeleteFileMeta(fileSha1)
}
