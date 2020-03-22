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

// Struct for user query parsing
type queryGenerator struct {
	once   sync.Once
	exp    *regexp.Regexp
	vector []string
}

func (qg *queryGenerator) remove3PSymbols(keys string) bool {
	var err error
	qg.once.Do(func() {
		qg.exp, err = regexp.Compile(`[,.;:\+\-&|~%@^$*(){}\[\]\\\/#<>"'` + "`" + `]`)
	})

	if err != nil {
		log.Println(err.Error())
		qg.vector = nil
		return false
	}

	keys = qg.exp.ReplaceAllString(keys, " ")
	qg.vector = strings.FieldsFunc(keys, func(r rune) bool {
		if r == ' ' {
			return true
		}
		return false
	})
	return true
}

func (qg *queryGenerator) getVector() []string {
	return qg.vector
}

func (qg *queryGenerator) generateArgsSql(itemNum int, operator string) string {
	valuesStr := ``
	for i := 1; i <= itemNum; i++ {
		valuesStr += `$` + strconv.Itoa(i) + ` `
		if i != itemNum {
			valuesStr += operator + ` `
		}
	}
	return valuesStr
}

// Generate args slice from words vector, using limit and page number,
// for reusing getEvents method from sqlEventsRepository struct
func (qg *queryGenerator) generateArgSlice(limit int, page int) []interface{} {
	keys := []string{
		strings.Join(qg.vector, " "),
		strconv.Itoa(limit),
		strconv.Itoa(limit * (page - 1)),
	}

	args := make([]interface{}, 3)
	for i, v := range keys {
		args[i] = v
	}
	return args
}

// +/- universal method for getting events array by condition (aka sqlStatement)
// and parameters in args (interface array)
func (er *sqlEventsRepository) getEvents(sqlStatement string, args ...interface{}) ([]models.Event, error) {
	baseSql := `SELECT eid, uid, title, edate, message, is_edited, author, etype, range, photos FROM events `
	baseSql += sqlStatement
	log.Println(baseSql, args)
	rows, err := er.db.Query(baseSql , args...)
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

// Getting all events without key words ordered by date
// Deprecated: DO NOT USE IN THE PRODUCTION MODE
func (er *sqlEventsRepository) GetAllEvents() ([]models.Event, error) {
	sqlCondition := ` ORDER BY edate ;`
	return er.getEvents(sqlCondition)
}

// Getting feed events. For now feed mean that this events ordered by date
// Thus it's closest events
func (er *sqlEventsRepository) GetFeedEvents(limit int, page int) ([]models.Event, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("Page number can't be less than 1\n")
	}

	sqlCondition := ` WHERE edate >= current_timestamp ORDER BY edate ASC LIMIT $1 OFFSET $2 ;`
	// TODO: add cool feed algorithm (aka select)
	return er.getEvents(sqlCondition, limit, page)
}

func (er *sqlEventsRepository) GetEventsByKeyWord(keyWordsString string, page int) ([]models.Event, error) {
	log.Println(keyWordsString)
	log.Println(page)
	if page < 1 {
		return nil, errors.New("Page number can't be less than 1\n")
	}

	sqlCondition := ` WHERE edate >= current_timestamp AND title_tsv @@ phraseto_tsquery('russian', $1 )
							ORDER BY edate ASC LIMIT $2 OFFSET $3 ;`

	var generator queryGenerator
	ok := generator.remove3PSymbols(keyWordsString)
	if !ok {
		return nil, errors.New("Incorrect symbols in the query\n")
	}

	args := generator.generateArgSlice(settings.UseCaseConf.PageLimit, page)
	return er.getEvents(sqlCondition, args...)
}
