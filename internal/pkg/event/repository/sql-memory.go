package repository

//go:generate mockgen -destination=../mocks/mock_repository.go -package=mocks failless/internal/pkg/event Repository

import (
	"errors"
	chatRepository "failless/internal/pkg/chat/repository"
	"failless/internal/pkg/event"
	"failless/internal/pkg/models"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	QueryInsertMidEventMember = `
		INSERT INTO		mid_event_members (uid, eid)
		VALUES			($1, $2)
		ON CONFLICT
		ON CONSTRAINT 	unique_member
		DO				NOTHING;`
	QueryIncrementMemberAmount = `
		UPDATE		mid_events
		SET			members = members + 1
		WHERE		eid = $1 `
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
		return r == ' '
	})
	return true
}

//func (qg *queryGenerator) getVector() []string {
//	return qg.vector
//}

//func (qg *queryGenerator) generateArgsSql(itemNum int, operator string, option string, offset int) string {
//	sql := ``
//	for i := 1 + offset; i <= itemNum; i++ {
//		sql += option + `$` + strconv.Itoa(i) + ` `
//		if i != itemNum {
//			sql += operator + ` `
//		}
//	}
//	return sql
//}

// Getting AND condition for time of event and sorting by date with pagination
//func (qg *queryGenerator) getConstantCondition(itemNum int) string {
//	return `
//		AND 		e.edate >= current_timestamp
//		ORDER BY 	e.edate ASC
//		LIMIT 		$` + strconv.Itoa(itemNum+1) + `
//		OFFSET 		$` + strconv.Itoa(itemNum+2) + ` ;`
//}

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

func (er *sqlEventsRepository) GetNameByID(uid int) (string, error) {
	sqlStatement := `
		SELECT 	name
		FROM 	profile
		WHERE 	uid = $1 ;`
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

func (er *sqlEventsRepository) GetValidTags() ([]models.Tag, error) {
	sqlStatement := `
		SELECT 		tag_id, name
		FROM 		tag
		ORDER BY 	tag_id;`
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

//func (er *sqlEventsRepository) generateORStatement(fieldName string, length int) string {
//	sql := ``
//	for i := 1; i <= length; i++ {
//		sql += fieldName + " = $" + strconv.Itoa(i)
//		if i != length {
//			sql += " OR "
//		}
//	}
//	return sql
//}

func (er *sqlEventsRepository) FollowBigEvent(uid, eid int) error {
	sqlStatement := `
		INSERT INTO subscribe (uid, table_id)
		VALUES 		( $1, $2 );`

	rows, err := er.db.Exec(sqlStatement, uid, eid)
	if err != nil || rows.RowsAffected() == 0 {
		return err
	}
	return nil
}

func (er *sqlEventsRepository) UnfollowMidEvent(uid, eid int) error {
	sqlStatement := `
		UPDATE		event_vote
		SET			value = -1, is_edited = TRUE, vote_date = current_timestamp
		WHERE		uid = $1
		AND			eid = $2;`
	rows, err := er.db.Exec(sqlStatement, uid, eid)
	if err != nil || rows.RowsAffected() == 0 {
		log.Println(err)
		log.Println(sqlStatement, uid, eid)
		return err
	}

	return nil
}

func (er *sqlEventsRepository) UnfollowBigEvent(uid, eid int) error {
	sqlStatement := `
		DELETE FROM		subscribe
		WHERE			uid = $1
		AND				table_id = $2;`

	rows, err := er.db.Exec(sqlStatement, uid, eid)
	if err != nil || rows.RowsAffected() == 0 {
		return err
	}
	return nil
}

func (er *sqlEventsRepository) CreateSmallEvent(event *models.SmallEvent) error {
	sqlStatement := `
		INSERT INTO		small_event (uid, title, description, date, tags, photos)
		VALUES			($1, $2, $3, $4, $5, $6)
		RETURNING		eid;`

	row := er.db.QueryRow(sqlStatement, event.UId, event.Title, event.Descr, event.Date, event.TagsId, event.Photos)
	err := row.Scan(
		&event.EId)
	if err != nil {
		log.Println(row)
		return err
	}

	return nil
}

func (er *sqlEventsRepository) UpdateSmallEvent(event *models.SmallEvent) (int, error) {
	sqlStatement := `
		UPDATE	small_event
		SET		title = $3, description = $4, date = $5, tags = $6, photos = $7
		WHERE	uid = $1
		AND		eid = $2;`

	cTag, err := er.db.Exec(sqlStatement, event.UId, event.EId, event.Title, event.Descr, event.Date, event.TagsId, event.Photos)
	if err != nil || cTag.RowsAffected() == 0 {
		log.Println(err)
		return http.StatusNotFound, err
	}

	// TODO: try it

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) DeleteSmallEvent(uid int, eid int64) error {
	sqlStatement := `
		DELETE FROM	small_event
		WHERE		uid = $1
		AND			eid = $2;`

	cTag, err := er.db.Exec(sqlStatement, uid, eid)
	if err != nil || cTag.RowsAffected() == 0 {
		log.Println(err)
		return err
	}

	// TODO: try it

	return nil
}

func (er *sqlEventsRepository) GetSmallEventsForUser(smallEvents *models.SmallEventList, uid int) (int, error) {
	sqlStatement := `
		SELECT		eid, uid, title, description, date, tags, photos
		FROM		small_event
		WHERE		uid = $1
		ORDER BY	(current_timestamp - date) ASC, time_created DESC
		LIMIT		30;`

	rows, err := er.db.Query(sqlStatement, uid)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	defer rows.Close()

	tags := pgtype.Int4Array{}
	date := new(time.Time)
	for rows.Next() {
		eventInfo := models.SmallEvent{}
		err = rows.Scan(
			&eventInfo.EId,
			&eventInfo.UId,
			&eventInfo.Title,
			&eventInfo.Descr,
			&date,
			&tags,
			&eventInfo.Photos)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, err
		}

		if date != nil {
			eventInfo.Date = *date
		}

		if err = tags.AssignTo(&eventInfo.TagsId); err != nil {
			log.Println(err)
		}
		*smallEvents = append(*smallEvents, eventInfo)
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) CreateMidEvent(event *models.MidEvent) error {
	tx, err := er.db.Begin()
	if err != nil {
		log.Println("CreateMidEvent: Failed to start transaction", err)
		return err
	}
	//defer tx.Rollback()
	defer func() {
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
	}()
	// Create global chat
	var chatID int64
	row := tx.QueryRow(chatRepository.QueryInsertGlobalChat, event.AdminId, event.Limit, event.Title)
	if err = row.Scan(&chatID); err != nil {
		log.Println("CreateMidEvent: ", err)
		log.Println(row)
		return err
	}

	// Create local chat
	var userLocalID int
	var photo *string
	if event.Photos != nil {
		photo = &event.Photos[0]
	} else {
		photo = nil
	}
	row = tx.QueryRow(chatRepository.QueryInsertNewLocalChat, chatID, event.AdminId, photo, event.Title)
	if err = row.Scan(&userLocalID); err != nil {
		log.Println("CreateMidEvent: ", err)
		log.Println(row)
		return err
	}

	//Insert first message
	if row, err_QueryInsertFirstMessage := tx.Exec(chatRepository.QueryInsertFirstMessage, event.AdminId, chatID, userLocalID, true); err != nil {
		log.Println("CreateMidEvent: ", err_QueryInsertFirstMessage)
		log.Println(row)
		return err_QueryInsertFirstMessage
	}

	// Create event
	sqlStatement := `
		INSERT INTO	mid_events (admin_id, title, description, date, tags, photos, member_limit, is_public, chat_id, title_tsv)
		SELECT		$1, $2, $3, $4, $5, $6, $7, $8, $9,
					setweight(to_tsvector($10), 'A')
						|| 
					setweight(to_tsvector($11), 'B')
		RETURNING	eid, members;`
	row = tx.QueryRow(sqlStatement,
		event.AdminId,
		event.Title,
		event.Descr,
		event.Date,
		event.TagsId,
		event.Photos,
		event.Limit,
		event.Public,
		chatID,
		event.Title,
		event.Descr)
	if err = row.Scan(&event.EId, &event.MemberAmount); err != nil {
		log.Println("CreateMidEvent: ", err)
		log.Println(row)
		return err
	}

	// Add first member
	if cTag, err_QueryInsertMidEventMember := tx.Exec(QueryInsertMidEventMember, event.AdminId, event.EId); err != nil || cTag.RowsAffected() == 0 {
		log.Println("CreateMidEvent: ", err_QueryInsertMidEventMember)
		return err_QueryInsertMidEventMember
	}

	// Update eid in global chats
	if cTag, err_QueryUpdateEidAvatarInChatUser := tx.Exec(chatRepository.QueryUpdateEidAvatarInChatUser, event.EId, photo, chatID); err != nil || cTag.RowsAffected() == 0 {
		log.Println("CreateMidEvent: ", err_QueryUpdateEidAvatarInChatUser)
		return err_QueryUpdateEidAvatarInChatUser
	}

	// Close transaction
	if err = tx.Commit(); err != nil {
		log.Println("CreateMidEvent: ", err)
		return err
	}

	return nil
}

func (er *sqlEventsRepository) GetOwnMidEvents(midEvents *models.MidEventList, uid int) (int, error) {
	sqlStatement := `
		SELECT		eid, title, description, tags, date, photos, member_limit, members, is_public
		FROM		mid_events
		WHERE		admin_id = $1
		ORDER BY	(current_timestamp - date) ASC, time_created DESC
		LIMIT		30;`
	rows, err := er.db.Query(sqlStatement, uid)
	if err != nil {
		log.Println("EventRepo: GetOwnMidEvents: ", err)
		return http.StatusInternalServerError, err
	}
	defer rows.Close()

	err = er.retrieveMidEventsFrom(rows, midEvents, "none")

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) GetAllMidEvents(midEvents *models.MidEventList, request *models.EventRequest) (int, error) {
	sqlStatement := `
		SELECT		eid, title, description, tags, date, photos, member_limit, members, is_public
		FROM		mid_events
		`
	var rows *pgx.Rows
	var err error
	if len(request.Query) > 0 {
		var generator queryGenerator
		if !generator.remove3PSymbols(request.Query) {
			return http.StatusInternalServerError, errors.New("Incorrect symbols in the query\n")
		}
		args := generator.GenerateArgSlice(request.Limit, request.Page)
		sqlStatement += `
		WHERE		e.title_tsv @@ phraseto_tsquery( $1 )
		ORDER BY	(current_timestamp - date) ASC, time_created DESC
		LIMIT		$2
		OFFSET		$3 * 0;` // TODO: offset 0
		rows, err = er.db.Query(sqlStatement, args...)
	} else {
		sqlStatement += `
		ORDER BY	(current_timestamp - date) ASC, time_created DESC
		LIMIT		$1
		OFFSET		$2 * 0;`	// TODO: offset 0
		rows, err = er.db.Query(sqlStatement, request.Limit, request.Page)
	}

	if err != nil {
		log.Println("EventRepo: GetAllMidEvents: ", err)
		return http.StatusInternalServerError, err
	}
	defer rows.Close()

	err = er.retrieveMidEventsFrom(rows, midEvents, "false")

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) GetSubscriptionMidEvents(midEvents *models.MidEventList, uid int) (int, error) {
	sqlStatement := `
		SELECT		mid_events.eid, title, description, tags, date, photos, member_limit, members, is_public
		FROM		mid_events
		JOIN		mid_event_members
		ON			mid_events.eid = mid_event_members.eid
		AND			mid_event_members.uid = $1
		ORDER BY	(current_timestamp - date) ASC, time_created DESC
		LIMIT		30;`
	rows, err := er.db.Query(sqlStatement, uid)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	defer rows.Close()

	err = er.retrieveMidEventsFrom(rows, midEvents, "true")

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) GetMidEventsWithFollowed(midEvents *models.MidEventList, request *models.EventRequest) (int, error) {
	sqlStatement := `
		SELECT		me.eid, me.title, me.description, me.tags, me.date, me.photos, me.member_limit, me.members, me.is_public,
					CASE WHEN me_mem.uid IS NULL THEN FALSE ELSE TRUE END AS followed
		FROM 		mid_events AS me
		LEFT JOIN 	mid_event_members AS me_mem
		ON			me_mem.eid = me.eid AND me_mem.uid = $1
		WHERE 		me.admin_id <> $1`
	var rows *pgx.Rows
	var err error
	if len(request.Query) > 0 {
		var generator queryGenerator
		if !generator.remove3PSymbols(request.Query) {
			return http.StatusInternalServerError, errors.New("Incorrect symbols in the query\n")
		}
		args := generator.GenerateArgSlice(request.Limit, request.Page)
		sqlStatement += `
		AND 		me.title_tsv @@ phraseto_tsquery( $2 )
		ORDER BY	(current_timestamp - me.date) ASC, me.time_created DESC
		LIMIT		$3
		OFFSET		$4 * 0;` // TODO: offset 0 (it's by default 0 .-.)
		args = append([]interface{}{request.Uid}, args...)
		rows, err = er.db.Query(sqlStatement, args...)
	} else {
		sqlStatement += `
		ORDER BY	(current_timestamp - me.date) ASC, me.time_created DESC
		LIMIT		$2 
		OFFSET		$3 * 0;` // TODO: offset 0
		rows, err = er.db.Query(sqlStatement, request.Uid, request.Limit, request.Page)
	}
	if err != nil {
		log.Println("EventRepo: GetMidEventsWithFollowed: ", err)
		log.Println(sqlStatement)
		return http.StatusInternalServerError, err
	}
	defer rows.Close()

	err = er.retrieveMidEventsFrom(rows, midEvents, "find")

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) JoinMidEvent(uid, eid int) (int, error) {
	tx, err := er.db.Begin()
	if err != nil {
		log.Println("JoinMidEvent: Failed to start transaction", err)
		return http.StatusInternalServerError, err
	}
	//defer tx.Rollback()
	defer func() {
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
	}()
	cTag, err := tx.Exec(QueryInsertMidEventMember, uid, eid)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	if cTag.RowsAffected() == 0 {
		log.Println("Member already in base", cTag)
	} else {
		var chatID int
		var eventTitle string
		sqlStatement := QueryIncrementMemberAmount + `
		RETURNING 	chat_id, title;`
		if err = er.db.QueryRow(sqlStatement, eid).Scan(&chatID, &eventTitle); err != nil {
			log.Println(err)
			return http.StatusInternalServerError, err
		}

		var userLocalID int
		row := tx.QueryRow(chatRepository.QueryInsertNewLocalChat, chatID, uid, nil, eventTitle)
		if err = row.Scan(&userLocalID); err != nil {
			log.Println("JoinMidEvent: ", err)
			log.Println(row)
			return http.StatusInternalServerError, err
		}

		//Insert first message
		if row, err_QueryInsertFirstMessage := tx.Exec(chatRepository.QueryInsertFirstMessage, uid, chatID, userLocalID, true); err != nil {
			log.Println("JoinMidEvent: ", err_QueryInsertFirstMessage)
			log.Println(row)
			return http.StatusInternalServerError, err_QueryInsertFirstMessage
		}
	}

	// Close transaction
	if err = tx.Commit(); err != nil {
		log.Println("JoinMidEvent: ", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) LeaveMidEvent(uid, eid int) (int, error) {
	sqlStatement := `
		DELETE FROM	mid_event_members
		WHERE		uid = $1
		AND			eid = $2;`
	cTag, err := er.db.Exec(sqlStatement, uid, eid)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	if cTag.RowsAffected() == 0 {
		log.Println("User wasn't a member of mid-event", cTag)
	} else {
		sqlStatement = `
		UPDATE	mid_events
		SET		members = members - 1
		WHERE	eid = $1;`
		cTag, err = er.db.Exec(sqlStatement, eid)
		if err != nil || cTag.RowsAffected() == 0 {
			log.Println(err)
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) GetOwnMidEventsWithAnotherUserFollowed(midEvents *models.MidEventList, admin, member int) (int, error) {
	sqlStatement := `
		SELECT		ME.eid, title, description, tags, date, photos, member_limit, members, is_public,
					CASE WHEN MEM.uid IS NULL THEN FALSE ELSE TRUE END AS followed
		FROM		mid_events ME
		LEFT JOIN	mid_event_members MEM
		ON			MEM.eid = ME.eid
		AND			MEM.uid = $2
		WHERE		admin_id = $1
		ORDER BY	(current_timestamp - date) ASC, time_created DESC
		LIMIT		20;`
	rows, err := er.db.Query(sqlStatement, admin, member)

	if err != nil {
		log.Println("GetOwnMidEventsWithAnotherUserFollowed", err)
		return http.StatusInternalServerError, err
	}
	defer rows.Close()

	err = er.retrieveMidEventsFrom(rows, midEvents, "find")
	if err != nil {
		log.Println("GetOwnMidEventsWithAnotherUserFollowed", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) GetSubscriptionMidEventsWithAnotherUserFollowed(midEvents *models.MidEventList, uid, visitor int) (int, error) {
	sqlStatement := `
		SELECT		grouped.eid, title, description, tags, date, photos, member_limit, members, is_public, followed
		FROM
			(SELECT		ME.eid, title, description, tags, date, photos, member_limit, members, is_public,
						CASE WHEN COUNT(ME.eid) = 1 THEN FALSE ELSE TRUE END AS followed
			FROM		mid_events ME
			JOIN		mid_event_members MEM
			ON			ME.eid = MEM.eid
			AND			(MEM.uid = $1
			OR			MEM.uid = $2)
			GROUP BY	ME.eid, date, time_created) AS grouped
		JOIN		mid_event_members MEM
		ON			MEM.eid = grouped.eid
		AND			MEM.uid = $1
		ORDER BY	(current_timestamp - date) ASC
		LIMIT		30;`
	rows, err := er.db.Query(sqlStatement, uid, visitor)
	if err != nil {
		log.Println("GetSubscriptionMidEventsWithAnotherUserFollowed", err)
		return http.StatusInternalServerError, err
	}
	defer rows.Close()

	err = er.retrieveMidEventsFrom(rows, midEvents, "find")

	if err != nil {
		log.Println("GetSubscriptionMidEventsWithAnotherUserFollowed", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (er *sqlEventsRepository) retrieveMidEventsFrom(rows *pgx.Rows, midEvents *models.MidEventList, followed string) error {
	var err error
	tags := pgtype.Int4Array{}
	date := new(time.Time)
	for rows.Next() {
		eventInfo := models.MidEvent{}
		if followed == "find" {
			err = rows.Scan(
				&eventInfo.EId,
				&eventInfo.Title,
				&eventInfo.Descr,
				&tags,
				&date,
				&eventInfo.Photos,
				&eventInfo.Limit,
				&eventInfo.MemberAmount,
				&eventInfo.Public,
				&eventInfo.Followed)
		} else {
			err = rows.Scan(
				&eventInfo.EId,
				&eventInfo.Title,
				&eventInfo.Descr,
				&tags,
				&date,
				&eventInfo.Photos,
				&eventInfo.Limit,
				&eventInfo.MemberAmount,
				&eventInfo.Public)
		}
		if err != nil {
			return err
		}

		if date != nil {
			eventInfo.Date = *date
		}

		if err = tags.AssignTo(&eventInfo.TagsId); err != nil {
			log.Println("EventRepo: GetMidEventsWithFollowed: Assigning tags:", err)
		}

		if followed == "true" {
			eventInfo.Followed = true
		} else if followed == "false" {
			eventInfo.Followed = false
		}

		*midEvents = append(*midEvents, eventInfo)
	}

	return nil
}
