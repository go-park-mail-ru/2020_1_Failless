package repository

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/user"
	"failless/internal/pkg/vote"
	"github.com/jackc/pgx"
	"log"
	"time"
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

func (vr *sqlVoteRepository) AddUserVote(uid int, id int, value int8) error {
	sqlStatement := `INSERT INTO user_vote (uid, user_id, value) VALUES ( $1 , $2 , $3 );`
	_, err := vr.db.Exec(sqlStatement, uid, id, value)
	if err != nil {
		log.Println(sqlStatement, uid, id, value)
		return err
	}

	return nil
}

func (vr *sqlVoteRepository) FindFollowers(eid int) ([]models.UserGeneral, error) {
	sqlStatement := `SELECT u.uid, u.name, p.about, p.gender, p.birthday, p.photos FROM event_vote AS ev
							NATURAL JOIN profile AS u JOIN profile_info AS p ON p.pid = u.uid 
							WHERE ev.eid = $1 AND ev.value > 0 ORDER BY ev.vote_date DESC;`

	rows, err := vr.db.Query(sqlStatement, eid)
	if err != nil && rows != nil && !rows.Next() {
		log.Println(sqlStatement)
		log.Println("event has not got any followers")
		return nil, nil
	} else if err != nil || rows == nil {
		log.Println(err.Error())
		return nil, err
	}

	var profiles []models.UserGeneral
	for rows.Next() {
		profile := models.UserGeneral{}
		gender := ""
		genderPtr := &gender
		bday := time.Time{}
		bdayPtr := &bday
		err = rows.Scan(
			&profile.Uid,
			&profile.Name,
			&profile.About,
			&genderPtr,
			&bdayPtr,
			&profile.Photos)
		if err != nil {
			log.Println("Error while getting profiles")
			return nil, err
		}
		if bdayPtr != nil {
			profile.Birthday = bday
		}
		if genderPtr != nil {
			profile.Gender = user.GenderByStr(gender)
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (vr *sqlVoteRepository) AddUserToChat(eid int, uid int) (models.Chat, error) {
	return models.Chat{}, nil
}
