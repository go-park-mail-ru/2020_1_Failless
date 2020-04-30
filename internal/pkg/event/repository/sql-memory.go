package repository

import (
	"errors"
	"failless/internal/pkg/event"
	"failless/internal/pkg/models"
	"failless/internal/pkg/settings"
	"fmt"
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

func (qg *queryGenerator) generateArgsSql(itemNum int, operator string, option string, offset int) string {
	sql := ``
	for i := 1 + offset; i <= itemNum; i++ {
		sql += option + `$` + strconv.Itoa(i) + ` `
		if i != itemNum {
			sql += operator + ` `
		}
	}
	return sql
}

// Getting AND condition for time of event and sorting by date with pagination
func (qg *queryGenerator) getConstantCondition(itemNum int) string {
	return `AND e.edate >= current_timestamp ORDER BY e.edate ASC LIMIT $` + strconv.Itoa(itemNum+1) + ` 
						OFFSET $` + strconv.Itoa(itemNum+2) + ` ;`
}

// Generate args slice from words vector, using limit and page number,
// for reusing getEvents method from sqlEventsRepository struct
func (qg *queryGenerator) GenerateArgSlice(limit int, page int) []interface{} {
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

func (qg *queryGenerator) JoinIntArgs(items []int, limit int, page int) []interface{} {
	keys := []int{
		limit,
		limit * (page - 1),
	}

	args := make([]interface{}, 2+len(items))
	for i, v := range items {
		args[i] = v
	}

	for i, v := range keys {
		args[i+len(items)] = v
	}
	return args
}

// +/- universal method for getting events array by condition (aka sqlStatement)
// and parameters in args (interface array)
func (er *sqlEventsRepository) getEvents(withCondition string, sqlStatement string, args ...interface{}) ([]models.Event, error) {
	baseSql := withCondition + ` SELECT e.eid, e.uid, e.title, e.edate, e.message, e.is_edited,
						e.author, e.etype, e.range, e.photos FROM events AS e `
	baseSql += sqlStatement
	log.Println(baseSql, args)
	rows, err := er.db.Query(baseSql, args...)
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
	sqlStatement := `INSERT INTO events (uid, title, message, author, etype, is_public, range, edate, title_tsv)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8,
							setweight(to_tsvector($9), 'A') || 
							setweight(to_tsvector($10), 'B')) RETURNING eid;`
	err := er.db.QueryRow(sqlStatement,
		event.AuthorId,
		event.Title,
		event.Message,
		event.Author,
		event.Type,
		event.Public,
		event.Limit,
		event.EDate,
		event.Title,
		event.Message).Scan(&event.EId)
	if err != nil {
		log.Println(err.Error())
		log.Println(sqlStatement, event.AuthorId, event.Title, event.Message,  event.Author,
			event.Type, event.Public,  event.Limit,  event.EDate)
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
	sqlCondition := ` ORDER BY e.edate ;`
	return er.getEvents("", sqlCondition)
}

// Getting vote events. For now vote mean that this events ordered by date
// Thus it's closest events
func (er *sqlEventsRepository) GetFeedEvents(uid int, limit int, page int) ([]models.Event, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("Page number can't be less than 1\n")
	}
	withCondition := `WITH voted_events AS ( SELECT eid FROM event_vote WHERE uid = $1 ) `
	sqlCondition := ` LEFT JOIN voted_events AS v ON e.eid = v.eid WHERE v.eid IS NULL AND e.uid != $2 AND
						e.edate >= current_timestamp ORDER BY e.edate ASC LIMIT $3 OFFSET $4 ;`
	// TODO: add cool vote algorithm (aka select)
	return er.getEvents(withCondition, sqlCondition, uid, uid, limit, page)
}

func (er *sqlEventsRepository) GetEventsByKeyWord(keyWords string, page int) (models.EventList, error) {
	log.Println(keyWords)
	log.Println(page)
	if page < 1 {
		return nil, errors.New("Page number can't be less than 1\n")
	}

	sqlCondition := ` WHERE e.edate >= current_timestamp AND e.title_tsv @@ phraseto_tsquery( $1 )
							ORDER BY e.edate ASC LIMIT $2 OFFSET $3 ;`

	var generator queryGenerator
	ok := generator.remove3PSymbols(keyWords)
	if !ok {
		return nil, errors.New("Incorrect symbols in the query\n")
	}

	args := generator.GenerateArgSlice(settings.UseCaseConf.PageLimit, page)
	return er.getEvents("", sqlCondition, args...)
}

func (er *sqlEventsRepository) GetValidTags() ([]models.Tag, error) {
	sqlStatement := `SELECT tag_id, name FROM tag ORDER BY tag_id;`
	rows, err := er.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	var tags []models.Tag
	for rows.Next() {
		tag := models.Tag{}
		err = rows.Scan(&tag.TagId, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (er *sqlEventsRepository) generateORStatement(fieldName string, length int) string {
	sql := ``
	for i := 1; i <= length; i++ {
		sql += fieldName + " = $" + strconv.Itoa(i)
		if i != length {
			sql += " OR "
		}
	}
	return sql
}

func (er *sqlEventsRepository) GetNewEventsByTags(tags []int, uid int, limit int, page int) (models.EventList, error) {
	var generator queryGenerator
	withCondition := `WITH voted_events AS ( SELECT eid FROM event_vote WHERE uid = $1 ) `
	sqlStatement := ` LEFT JOIN voted_events AS v ON e.eid = v.eid WHERE e.uid != $2 AND `
	items := append([]int{uid, uid}, tags...)
	sqlStatement += generator.generateArgsSql(len(items), "OR", "e.etype =", 2)
	sqlStatement += generator.getConstantCondition(len(items))
	return er.getEvents(withCondition, sqlStatement, generator.JoinIntArgs(items, limit, page)...)
}

func (er *sqlEventsRepository) FollowMidEvent(uid, eid int) error {
	sqlStatement := `INSERT INTO event_vote (uid, eid, value) VALUES ( $1 , $2 , 1 );`
	rows, err := er.db.Exec(sqlStatement, uid, eid)
	if err != nil || rows.RowsAffected() == 0 {
		log.Println(err)
		log.Println(sqlStatement, uid, eid)
		return err
	}

	return nil
}
func (er *sqlEventsRepository) FollowBigEvent(uid, eid int) error {
	sqlStatement := `
		INSERT INTO subscribe (uid, table_id) VALUES ( $1, $2 );`

	rows, err := er.db.Exec(sqlStatement, uid, eid)
	if err != nil || rows.RowsAffected() == 0 {
		return err
	}
	return nil
}

func (er *sqlEventsRepository) GetEventsWithFollowed(events *models.EventResponseList, request *models.EventRequest) error {
	sqlStatement := `
		SELECT e.eid, e.uid, e.title, e.edate, e.message, e.is_edited, e.author, e.etype, e.range, e.photos,
			   CASE WHEN ev.uid IS NULL THEN FALSE ELSE TRUE END AS followed
		FROM events e
		LEFT JOIN event_vote ev ON e.eid = ev.eid AND ev.uid = $1
		`

	var rows *pgx.Rows
	var err error
	if len(request.Query) > 0 {
		var generator queryGenerator
		if !generator.remove3PSymbols(request.Query) {
			return errors.New("Incorrect symbols in the query\n")
		}
		args := generator.GenerateArgSlice(request.Limit, request.Page)
		sqlStatement += "WHERE e.title_tsv @@ phraseto_tsquery( $2 ) LIMIT $3 OFFSET $4;"
		args = append([]interface{}{request.Uid}, args...)
		rows, err = er.db.Query(sqlStatement, args...)
	} else {
		sqlStatement += "LIMIT $2 OFFSET $3;"
		rows, err = er.db.Query(sqlStatement, request.Uid, request.Limit, request.Page)
	}
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return err
	}

	for rows.Next() {
		tempEvent := models.EventResponse{}
		err = rows.Scan(
			&tempEvent.Event.EId,
			&tempEvent.Event.AuthorId,
			&tempEvent.Event.Title,
			&tempEvent.Event.EDate,
			&tempEvent.Event.Message,
			&tempEvent.Event.Edited,
			&tempEvent.Event.Author,
			&tempEvent.Event.Type,
			&tempEvent.Event.Limit,
			&tempEvent.Event.Photos,
			&tempEvent.Followed)
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			return err
		}
		*events = append(*events, tempEvent)
	}

	return nil
}
