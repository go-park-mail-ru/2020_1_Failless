package repository

import (
	"failless/internal/pkg/event"
	"failless/internal/pkg/models"
	"github.com/jackc/pgx"
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
