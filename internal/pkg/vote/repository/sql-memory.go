package repository

import (
	"failless/internal/pkg/vote"
	"github.com/jackc/pgx"
	"log"
)

type sqlVoteRepository struct {
	db *pgx.ConnPool
}

func NewSqlVoteRepository(db *pgx.ConnPool) vote.Repository {
	return &sqlVoteRepository{db: db}
}

func (vr *sqlVoteRepository) AddUserVote(uid int, id int, value int8) error {
	sqlStatement := `INSERT INTO user_vote (uid, user_id, value) VALUES ( $1 , $2 , $3 );`
	_, err := vr.db.Exec(sqlStatement, uid, id, value)
	if err != nil {
		log.Println(sqlStatement, uid, id, value)
		return err
	}

	return nil
}

func (vr *sqlVoteRepository) CheckMatching(uid, id int) (bool, error) {
	sqlStatement := `
		SELECT uid, user_id
		FROM user_vote
		WHERE uid = $1 AND user_id = $2 AND value = 1;`

	rows, _ := vr.db.Exec(sqlStatement, id, uid)

	if rows.RowsAffected() == 1 {
		return true, nil
	}

	return false, nil
}
