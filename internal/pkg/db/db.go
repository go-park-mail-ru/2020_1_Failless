package db

import (
	"failless/internal/pkg/settings"
	"github.com/jackc/pgx"
	"log"
	"os"
	"sync"
)

var (
	// TODO: add it as environment variable
	host     = "pgbouncer"
	port     = uint16(6432)
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

var db *pgx.ConnPool = nil
var syncOnce = sync.Once{}

func init() {
	pgxConfig := pgx.ConnConfig{
		Host:     host,
		Port:     port,
		Database: dbname,
		User:     user,
		Password: password,
	}
	pgxConnPoolConfig := pgx.ConnPoolConfig{
		MaxConnections: 100,
		ConnConfig: pgxConfig,
	}
	dbase, err := pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		if settings.UseCaseConf.InHDD {
			log.Fatal("Connection to database was failed", pgxConfig)
		}
		db = nil
	} else {
		db = dbase
	}
}
func ConnectToDB() *pgx.ConnPool {
	if db == nil {
		log.Fatal("Connection to database was failed", host, port, user)
	}
	//syncOnce.Do(func() {
	//	pgxConfig := pgx.ConnConfig{
	//		Host:     host,
	//		Port:     port,
	//		Database: dbname,
	//		User:     user,
	//		Password: password,
	//	}
	//	pgxConnPoolConfig := pgx.ConnPoolConfig{
	//		MaxConnections: 100,
	//		ConnConfig: pgxConfig,
	//	}
	//	dbase, err := pgx.NewConnPool(pgxConnPoolConfig)
	//	if err != nil {
	//		if settings.UseCaseConf.InHDD {
	//			log.Fatal("Connection to database was failed")
	//		}
	//		db = nil
	//	} else {
	//		db = dbase
	//	}
	//})
	return db
}
