package meta

import (
	"fmt"
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
func UpdateFileMeta(meta FileMeta) {
	fileMetas[meta.FileSha1] = meta
}
func GetFileMeta(fileSha1 string) FileMeta {
	if meta, ok := fileMetas[fileSha1]; ok {
		fmt.Println(meta)
		return meta
	}
	return FileMeta{}
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
	return metaArray
}
