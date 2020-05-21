package repository

import (
	mydb "failless/internal/pkg/db"
	"failless/internal/pkg/vote"
	"log"
)

const (
	QueryInsertUserVote = `
		INSERT INTO user_vote (uid, user_id, value)
		VALUES 		( $1 , $2 , $3 );`
)

type sqlVoteRepository struct {
	db mydb.MyDBInterface
}

func NewSqlVoteRepository() vote.Repository {
	return &sqlVoteRepository{db: mydb.NewDBInterface()}
}

func (vr *sqlVoteRepository) AddUserVote(uid int, id int, value int8) error {
	_, err := vr.db.Exec(QueryInsertUserVote, uid, id, value)
	if err != nil {
		log.Println(err)
		log.Println(QueryInsertUserVote, uid, id, value)
		return err
	}

	return nil
}

func (vr *sqlVoteRepository) CheckMatching(uid, id int) (bool, error) {
	sqlStatement := `
		SELECT 	uid, user_id
		FROM 	user_vote
		WHERE 	uid = $1
		AND 	user_id = $2
		AND 	value = 1;`

	rows, _ := vr.db.Exec(sqlStatement, id, uid)

	if rows.RowsAffected() == 1 {
		return true, nil
	}

	return false, nil
}
