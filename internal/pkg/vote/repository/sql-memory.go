package repository

import (
	"failless/internal/pkg/models"
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

func (vr *sqlVoteRepository) AddEventVote(uid int, eid int, value int8) error {
	sqlStatement := `INSERT INTO event_vote (uid, eid, value) VALUES ( $1 , $2 , $3 );`
	_, err := vr.db.Exec(sqlStatement, uid, eid, value)
	if err != nil {
		log.Println(sqlStatement, uid, eid, value)
		return err
	}

	return nil
}

func (vr *sqlVoteRepository) FindFollowers(eid int) ([]models.User, error) {
	sqlStatement := `SELECT u.uid, u.name, u.phone, u.email FROM event_vote AS ev
							NATURAL JOIN profile AS u WHERE ev.eid = $1 AND ev.value > 0 ORDER BY ev.vote_date ASC;`

	rows, err := vr.db.Query(sqlStatement, eid)
	if err != nil && rows != nil && !rows.Next() {
		log.Println(sqlStatement)
		log.Println("event has not got any followers")
		return nil, nil
	} else if err != nil || rows == nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Uid, &user.Name, &user.Phone, &user.Email)
		if err != nil {
			log.Println("Error while getting users")
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (vr *sqlVoteRepository) AddUserToChat(eid int, uid int) (models.Chat, error) {
	return models.Chat{}, nil
}