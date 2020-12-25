package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (s *Sha1Stream) Update(data []byte) {
	if s._sha1 == nil {
		s._sha1 = sha1.New()
	}
	s._sha1.Write(data)
}
func (s *Sha1Stream) Sum() string {
	return hex.EncodeToString(s._sha1.Sum([]byte("")))
}
func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	_, _ = io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}
func FileMD5(file *os.File) string {
	_md5 := md5.New()
	_, _ = io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}
func FileSizeConversion(filesize int) string {
	if filesize/1024/1024/1024 >= 1 {
		return strconv.Itoa(filesize/1024/1024/1024) + "G"
	} else if filesize/1024/1024 >= 1 {
		return strconv.Itoa(filesize/1024/1024) + "M"
	} else if filesize/1024 >= 1 {
		return strconv.Itoa(filesize/1024) + "K"
	}
	return strconv.Itoa(filesize) + "B"
}
func GetFileSize(filename string) (result int64) {
	_ = filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		result = info.Size()
		return nil
	})
	return
}
