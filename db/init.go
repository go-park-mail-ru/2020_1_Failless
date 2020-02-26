package db

import (
	"github.com/jackc/pgx"
	"log"
	"sync"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "eventum"
	password = "eventum"
	dbname   = "eventum"
)
var db *pgx.ConnPool = nil
var syncOnce = sync.Once{}

func ConnectToDB() *pgx.ConnPool {
	syncOnce.Do(func() {
		pgxConfig := pgx.ConnConfig{
			Host:     host,
			Port:     port,
			Database: dbname,
			User:     user,
			Password: password,
		}
		pgxConnPoolConfig := pgx.ConnPoolConfig{
			ConnConfig: pgxConfig,
		}
		dbase, err := pgx.NewConnPool(pgxConnPoolConfig)
		if err != nil {
			log.Fatal("Connection to database was failed")
		}
		db = dbase
	})
	return db
}
