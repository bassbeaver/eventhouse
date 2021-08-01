package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	MysqlDbServiceAlias = "MysqlDbService"
)

func NewDB(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn) //"username:password@tcp(127.0.0.1:3306)/test"
	if err != nil {
		panic("Failed to connect to MySQL: " + err.Error())
	}

	return db
}
