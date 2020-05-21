package db

import (
	"github.com/jackc/pgx"
	"testing"
)

type MyDB struct {
	pgPool *pgx.ConnPool
}

func NewDBInterface() MyDBInterface {
	return &MyDB{pgPool:ConnectToDB()}
}

func (dtbs *MyDB) Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error) {
	return dtbs.pgPool.Exec(sql, arguments...)
}

var testConfig = pgx.ConnConfig {
	Host: "localhost",
	Port: uint16(5432),
	Database: "eventum_test_db",
}

var testConnPoolConfig = pgx.ConnPoolConfig{
	ConnConfig:     testConfig,
	MaxConnections: 1,
	AfterConnect:   nil,
	AcquireTimeout: 0,
}

func ConnectToTestDB(t *testing.T) *pgx.ConnPool {
	syncOnce.Do(func() {
		dbase, err := pgx.NewConnPool(testConnPoolConfig)
		if err != nil {
			t.Fatalf("Connection to database has failed %v", err)
		}
		db = dbase
	})

	return db
}

func CastPoolToConnPool(conn *pgx.Conn) *pgx.ConnPool {
	return nil
}
