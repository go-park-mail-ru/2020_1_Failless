package repository

import (
	"errors"
	"failless/internal/pkg/event"
	"failless/internal/pkg/models"
	"failless/internal/pkg/settings"
	"github.com/jackc/pgx"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type sqlEventsRepository struct {
	db *pgx.ConnPool
}

func NewSqlEventRepository(db *pgx.ConnPool) event.Repository {
	return &sqlEventsRepository{db: db}
}

func (er *sqlEventsRepository) getEvents(sqlStatement string, args ...interface{}) ([]models.Event, error) {
	log.Println(sqlStatement, args)
	rows, err := er.db.Query(sqlStatement, args...)
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
			&eventInfo.Limit,
			&eventInfo.Photos)
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

	return *namePtr, nil
}

// Get all events without key words ordered by date
func (er *sqlEventsRepository) GetAllEvents() ([]models.Event, error) {
	sqlStatement := `SELECT eid, uid, title, edate, message, is_edited, author, etype, range, photos FROM events ORDER BY edate ;`
	return er.getEvents(sqlStatement)
}

// Struct for user query parsing
type queryGenerator struct {
	once sync.Once
	exp  *regexp.Regexp
}

func (qg *queryGenerator) genAndQuery(keys string) ([]string, bool) {
	var err error
	qg.once.Do(func() {
		qg.exp, err = regexp.Compile(`[,.;:\+\-&|~%@^$*(){}\[\]\\\/#<>"'` + "`" + `]`)
	})

	if err != nil {
		log.Println(err.Error())
		return nil, false
	}

	keys = qg.exp.ReplaceAllString(keys, " ")
	result := strings.FieldsFunc(keys, func(r rune) bool {
		if r == ' ' {
			return true
		}
		return false
	})
	return result, true
}

func (qg *queryGenerator) generateSql(itemNum int, operator string) string {
	valuesStr := ``
	for i := 1; i <= itemNum; i++ {
		valuesStr += `$` + strconv.Itoa(i) + ` `
		if i != itemNum {
			valuesStr += operator + ` `
		}
	}
	return valuesStr
}

func (er *sqlEventsRepository) GetEventsByKeyWord(keyWordsString string, page int) ([]models.Event, error) {
	log.Println(keyWordsString)
	log.Println(page)
	if page < 1 {
		return nil, errors.New("Page number can't be less than 1\n")
	}

	// TODO: check for sql injections
	sqlStatement := `SELECT eid, uid, title, edate, message, is_edited, author, etype, range, photos FROM events
							WHERE edate >= current_timestamp `

	var keys []string
	clean := ""
	if keyWordsString != "" {
		sqlStatement += ` AND title_tsv @@ phraseto_tsquery('russian', $1 ) ORDER BY edate ASC LIMIT $2 OFFSET $3 ;`
		var generator queryGenerator
		vector, ok := generator.genAndQuery(keyWordsString)
		if !ok {
			return nil, errors.New("Incorrect symbols in the query\n")
		}
		clean = strings.Join(vector, " ")
		keys = append(keys, clean)
	} else {
		sqlStatement += ` ORDER BY edate ASC LIMIT $1 OFFSET $2 ;`
	}

	keys = append(keys, strconv.Itoa(settings.UseCaseConf.PageLimit))
	keys = append(keys, strconv.Itoa(settings.UseCaseConf.PageLimit*(page-1)))
	args := make([]interface{}, len(keys))
	for i, v := range keys {
		args[i] = v
	}

	return er.getEvents(sqlStatement, args...)
}
