package db

import (
	"github.com/aurora/Filestore-server/db/mysql"
	"github.com/aurora/Filestore-server/utils"
	"log"
)

func UserSignUp(username, password string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user(user_name,user_pwd) values(?,?)",
	)
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		log.Println(err)
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}
func UserSignIn(username, encPassword string) (ok bool, token string) {
	stmt, err := mysql.DBConn().Prepare(
		"select * from tbl_user where user_name=?",
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err)
		return
	} else if rows == nil {
		return
	}
	parseRows := mysql.ParseRows(rows)
	if len(parseRows) > 0 && string(parseRows[0]["user_pwd"].([]byte)) == encPassword {
		token, err = utils.GenerateToken(username, encPassword)
		if err != nil {
			log.Println(err)
			return
		}
		return true, token
	}
	return
}
