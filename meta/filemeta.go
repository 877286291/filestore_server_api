package meta

import (
	"github.com/aurora/Filestore-server/db"
	"time"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string `json:"-"`
	UploadAt string
}

func UpdateFileMetaDB(meta FileMeta) bool {
	return db.FileUploadFinished(meta.FileSha1, meta.FileName, meta.FileSize, meta.Location)
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
func RemoveFileMetaDB(fileSha1 string) bool {
	return db.DeleteFileMeta(fileSha1)
}
