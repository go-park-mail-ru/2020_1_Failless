package db

import "github.com/jackc/pgx"

type MyDBInterface interface {
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}
