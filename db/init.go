package db

import (
	"github.com/jackc/pgx"
	"log"
	"os"
	"sync"
)

var (
	host     = "localhost"
	port     = uint16(5432)
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
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
