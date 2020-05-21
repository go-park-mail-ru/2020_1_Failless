package db

//go:generate mockgen -destination=mocks/mock_db.go -package=mocks failless/internal/pkg/db MyDBInterface

import "github.com/jackc/pgx"

type MyDBInterface interface {
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}
