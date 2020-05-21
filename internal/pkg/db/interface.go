package db

//go:generate mockgen -destination=mocks/mock_db.go -package=mocks failless/internal/pkg/db MyDBInterface

import "github.com/jackc/pgx"

type MyDBInterface interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}
