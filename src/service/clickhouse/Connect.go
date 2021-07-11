package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
)

const (
	ClickhouseConnectServiceAlias = "ClickhouseConnect"
)

func NewClickhouseConnect(dsn string) *sql.DB {
	connect, err := sql.Open("clickhouse", dsn)
	if err != nil {
		panic(err)
	}

	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			panic(fmt.Sprintf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace))
		} else {
			panic(err)
		}
	}

	fmt.Println("Clickhouse connect created")

	return connect
}
