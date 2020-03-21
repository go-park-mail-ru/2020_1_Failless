package repository

import (
	"errors"
	"failless/internal/pkg/event"
	"failless/internal/pkg/models"
	"github.com/jackc/pgx"
	"log"
)

type sqlEventsRepository struct {
	db *pgx.ConnPool
}

func NewSqlEventRepository(db *pgx.ConnPool) event.Repository {
	return &sqlEventsRepository{db: db}
}

func (er *sqlEventsRepository) GetAllEvents() ([]models.Event, error) {
	sqlStatement := `SELECT eid, uid, title, edate, message, is_edited, author, etype, range FROM events ORDER BY edate ;`
	rows, err := er.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	var events []models.Event
	for rows.Next() {
		eventInfo := models.Event{}
		err = rows.Scan(
			&eventInfo.EId,
			&eventInfo.AuthorId,
			&eventInfo.Title,
			&eventInfo.EDate,
			&eventInfo.Message,
			&eventInfo.Edited,
			&eventInfo.Author,
			&eventInfo.Type,
			&eventInfo.Limit)
		if err != nil {
			return nil, err
		}
		events = append(events, eventInfo)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (er *sqlEventsRepository) SaveNewEvent(event *models.Event) error {
	sqlStatement := `INSERT INTO events (uid, title, message, author, etype, is_public, range)
							VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING eid, edate;`
	err := er.db.QueryRow(sqlStatement,
		event.AuthorId,
		event.Title,
		event.Message,
		event.Author,
		event.Type,
		event.Public,
		event.Limit).Scan(&event.EId, &event.EDate)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (er *sqlEventsRepository) GetNameByID(uid int) (string, error) {
	sqlStatement := `SELECT name FROM profile WHERE uid = $1 ;`
	var name string
	namePtr := &name
	err := er.db.QueryRow(sqlStatement, uid).Scan(&namePtr)
	if err != nil {
		log.Println(err.Error())
		return "", err
	} else if namePtr == nil {
		return "", errors.New("User not found\n")
	}

	//TODO: check is it works
	return name, nil
}
