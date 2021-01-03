package db

import (
	"github.com/aurora/Filestore-server/db/mysql"
	"log"
	"time"
)

type UserFile struct {
	UserName    string
	FileSha1    string
	FileSize    int64
	FileName    string
	UpdateAt    string
	LastUpdated string
}

func UserFileUploaded(username, fileSha1, filename string, filesize int64) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user_file(`user_name`,`file_sha1`,`file_name`,`file_size`,`upload_at`) values (?,?,?,?,?)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, fileSha1, filename, filesize, time.Now())
	if err != nil {
		return false
	}
	return true
}
func GetUserFileMetas(username string, limit int) (userFiles []UserFile, err error) {
	stmt, err := mysql.DBConn().Prepare(
		"select user_name, file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? order by upload_at desc limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		file := UserFile{}
		err = rows.Scan(&file.UserName, &file.FileSha1, &file.FileName, &file.FileSize, &file.UpdateAt, &file.LastUpdated)
		if err != nil {
			log.Println(err)
			break
		}
		userFiles = append(userFiles, file)
	}
	return
}
