package db

import "testing"

func TestConnectToDB(t *testing.T) {
	db := ConnectToDB()
	if db == nil {
		t.Fatal("Connection to database fail")
	}
}

