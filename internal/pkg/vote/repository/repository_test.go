package repository

import (
	"testing"
)

func TestSqlVoteRepository_AddEventVote(t *testing.T) {
	//testDB, mock, err := sqlmock.New()
	//if err != nil {
	//	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	//}
	//defer testDB.Close()
	//
	//conn, err := stdlib.AcquireConn(testDB)
	//if err != nil {
	//	t.Errorf("error was not expected durin acquiring connection: %v", err)
	//}
	//defer stdlib.ReleaseConn(testDB, conn)
	//
	//var vr = NewSqlVoteRepository(db.CastPoolToConnPool(conn))
	//
	//mock.ExpectBegin()
	//mock.ExpectExec("INSERT INTO event_vote (uid, eid, value)").WithArgs(1, 2, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	//mock.ExpectCommit()
	//
	//// now we execute our method
	//if err = vr.AddEventVote(1, 2, 1); err != nil {
	//	t.Errorf("error was not expected while updating stats: %s", err)
	//}
	//
	//// we make sure that all expectations were met
	//if err := mock.ExpectationsWereMet(); err != nil {
	//	t.Errorf("there were unfulfilled expectations: %s", err)
	//}
}

func TestSqlVoteRepository_FindFollowers(t *testing.T) {

}

func TestSqlVoteRepository_FindFollowers2(t *testing.T) {

}
