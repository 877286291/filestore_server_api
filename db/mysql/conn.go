package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fileserver?charset=utf8&parseTime=true")
	db.SetMaxOpenConns(10)
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
}
func DBConn() *sql.DB {
	return db
}
