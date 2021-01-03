package db

import (
	"database/sql"
	"github.com/aurora/Filestore-server/db/mysql"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	CreateAt sql.NullTime
	UpdateAt sql.NullTime
}

func GetFileMeta(fileHash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1,file_name,file_size,file_addr,create_at,update_at from tbl_file where file_sha1=? and status=1 limit 1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(fileHash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr, &tfile.CreateAt, &tfile.UpdateAt)
	if err != nil {
		return nil, err
	}
	return &tfile, nil
}
func FileUploadFinished(fileHash, filename string, filesize int64, fileAddr string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) value (?,?,?,?,1)",
	)
	if err != nil {
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileHash, filename, filesize, fileAddr)
	if err != nil {
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil && rf >= 0 {
		return true
	}
	return false
}
func DeleteFileMeta(fileHash string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"delete from tbl_file where file_sha1=?",
	)
	if err != nil {
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileHash)
	if err != nil {
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil && rf >= 0 {
		return true
	}
	return false
}
