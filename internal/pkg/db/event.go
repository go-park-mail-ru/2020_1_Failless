package db

import "github.com/jackc/pgx"

func GetAllEvents(db *pgx.ConnPool) ([]Event, error) {
	sqlStatement := `SELECT eid, uid, title, edate, message, is_edited, author, etype, range FROM events ORDER BY edate ;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	var events []Event
	for rows.Next() {
		event := Event{}
		err = rows.Scan(
			&event.EId,
			&event.AuthorId,
			&event.Title,
			&event.EDate,
			&event.Message,
			&event.Edited,
			&event.Author,
			&event.Type,
			&event.Limit)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return events, nil
}

