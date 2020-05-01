package repository

import (
	"errors"
	"failless/internal/pkg/models"
	"failless/internal/pkg/user"
	"github.com/jackc/pgx"
	"log"
	"time"
)

type sqlUserRepository struct {
	db *pgx.ConnPool
}

func NewSqlUserRepository(db *pgx.ConnPool) user.Repository {
	return &sqlUserRepository{db: db}
}

func (ur *sqlUserRepository) GetUserByUID(uid int) (models.User, error) {
	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE uid = $1;`
	return ur.getUser(sqlStatement, uid)
}

func (ur *sqlUserRepository) GetUserByPhoneOrEmail(phone string, email string) (models.User, error) {
	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE phone = $1 OR email = LOWER($2);`
	return ur.getUser(sqlStatement, phone, email)
}

// Private method
func (ur *sqlUserRepository) getUser(sqlStatement string, args ...interface{}) (user models.User, err error) {
	row := ur.db.QueryRow(sqlStatement, args...)
	err = row.Scan(
		&user.Uid,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Password)

	if err == pgx.ErrNoRows {
		return models.User{Uid: -1, Password: []byte{}}, nil
	} else if err != nil {
		return models.User{}, err
	}
	log.Println(user)
	return user, nil
}

// Private method
func (ur *sqlUserRepository) getEvents(sqlStatement string, args ...interface{}) ([]models.Event, error) {
	rows, err := ur.db.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	if rows != nil && !rows.Next() {
		rows.Close()
		log.Println(sqlStatement)
		log.Println("user has no events")
		return nil, nil
	} else if rows == nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		event := models.Event{}
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
			log.Println("Error while getting events")
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (ur *sqlUserRepository) AddNewUser(user *models.User) error {
	uid := 0
	sqlStatement := `INSERT INTO profile VALUES (default, $1, $2, LOWER($3), $4) RETURNING uid;`
	err := ur.db.QueryRow(sqlStatement, user.Name, user.Phone, user.Email, user.Password).Scan(&uid)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	user.Uid = uid
	sqlStatement = `INSERT INTO profile_info VALUES ( $1 , '' , default , default , default , default , default , default ) ;`
	_, err = ur.db.Exec(sqlStatement, user.Uid)
	if err != nil {
		log.Println(sqlStatement, user.Uid)
		return err
	}
	return nil
}

func (ur *sqlUserRepository) AddUserInfo(credentials models.User, info models.JsonInfo) error {
	if credentials.Uid < 0 {
		sqlStatement := `SELECT uid FROM profile WHERE LOWER(email) = LOWER($1) OR phone = $2;`
		row := ur.db.QueryRow(sqlStatement, credentials.Email, credentials.Phone)
		err := row.Scan(&credentials.Uid)
		if err == pgx.ErrNoRows {
			return errors.New("User " + credentials.Email + "doesn't exist")
		} else if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	gender := user.GenderById(info.Gender)

	sqlStatement := `UPDATE profile_info SET about = $1, photos = array_append(photos, $2), birthday = $3, gender = $4 WHERE pid = $5;`
	_, err := ur.db.Exec(sqlStatement, info.About, info.Photos[0], info.Birthday, gender, credentials.Uid)
	if err != nil {
		log.Println(sqlStatement, info.About, info.Photos[0], info.Birthday, info.Gender, credentials.Uid)
		log.Println(err.Error())
		return err
	}
	return nil
}

func (ur *sqlUserRepository) SetUserLocation(uid int, point models.LocationPoint) error {
	sqlStatement := `UPDATE profile_info SET location = ST_POINT($1, $2) WHERE pid = $3;`
	_, err := ur.db.Exec(sqlStatement, point.Latitude, point.Longitude, uid)
	return err
}

func (ur *sqlUserRepository) UpdateUserRating(uid int, rating float32) error {
	sqlStatement := `UPDATE profile_info SET rating = $1 WHERE pid = $2;`
	_, err := ur.db.Exec(sqlStatement, rating, uid)
	return err
}

func (ur *sqlUserRepository) UpdateUserTags(uid int, tagId int) error {
	sqlStatement := `CALL update_tag($1, $2);`
	_, err := ur.db.Exec(sqlStatement, uid, tagId)
	return err
}

func (ur *sqlUserRepository) UpdateUserSimple(uid int, social []string, about *string) error {
	sqlStatement := `UPDATE profile_info SET about = $1, social = $2 WHERE pid = $3;`
	_, err := ur.db.Exec(sqlStatement, *about, social, uid)
	return err
}

func (ur *sqlUserRepository) GetProfileInfo(uid int) (info models.JsonInfo, err error) {
	var profile ProfileInfo
	sqlStatement := `SELECT about, photos, rating, birthday, gender FROM profile_info WHERE pid = $1 ;`
	gender := ""
	genPtr := &gender
	err = ur.db.QueryRow(sqlStatement, uid).Scan(
		&profile.About,
		&profile.Photos,
		&profile.Rating,
		&profile.Birthday,
		&genPtr)
	if err != nil {
		log.Println(sqlStatement)
		log.Println("error in get profile info")
		return models.JsonInfo{}, err
	}

	if profile.About != nil {
		info.About = *profile.About
	}

	if profile.Photos != nil {
		info.Photos = *profile.Photos
	}

	if genPtr != nil {
		info.Gender = user.GenderByStr(gender)
	}

	if profile.Location != nil {
		info.Location = *profile.Location
	}
	return
}

// func for testing
func (ur *sqlUserRepository) DeleteUser(mail string) error {
	sqlStatement := `DELETE FROM profile WHERE email=$1;`
	_, err := ur.db.Exec(sqlStatement, mail)
	return err
}

// TODO: move it to event pkg
func (ur *sqlUserRepository) GetUserEvents(uid int) ([]models.Event, error) {
	sqlStatement := `SELECT eid, uid, title, edate, message, is_edited, author, etype, range FROM events WHERE uid = $1 LIMIT 10;`
	return ur.getEvents(sqlStatement, uid)
}

// TODO: move it to event pkg
func (ur *sqlUserRepository) GetEventsByTag(tag string) ([]models.Event, error) {
	sqlStatement := `SELECT eid, uid, title, edate, message, is_edited, author, etype, range FROM events WHERE etype = $1 ;`
	return ur.getEvents(sqlStatement, tag)
}

func (ur *sqlUserRepository) GetUserTags(uid int) ([]models.Tag, error) {
	sqlStatement := `SELECT tag_id, name FROM user_tag NATURAL JOIN tag WHERE uid = $1 ;`
	rows, err := ur.db.Query(sqlStatement, uid)
	if err != nil {
		return nil, err
	}
	if rows != nil && !rows.Next() {
		log.Println(sqlStatement)
		log.Println("user has no tags")
		return nil, nil
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		tag := models.Tag{}
		err = rows.Scan(&tag.TagId, &tag.Name)
		if err != nil {
			log.Println("Error while getting tags")
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (ur *sqlUserRepository) UpdateUserPhotos(uid int, name string) error {
	sqlStatement := `UPDATE profile_info SET photos = array_append(photos, $1) WHERE pid = $2;`
	_, err := ur.db.Exec(sqlStatement, name, uid)
	if err != nil {
		log.Println(sqlStatement, name, uid)
		log.Println(err.Error())
		return err
	}
	return nil
}

func (ur *sqlUserRepository) UpdUserGeneral(info models.JsonInfo, usr models.User) error {
	//gender := user.GenderById(info.Gender)

	tx, err := ur.db.Begin()
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()

	// TODO: add name
	sqlStatement := `UPDATE profile SET email = $1, phone = $2, password = $3 WHERE uid = $4;`
	_, err = tx.Exec(sqlStatement, usr.Email, usr.Phone, usr.Password)
	if err != nil {
		log.Println(sqlStatement, usr.Email, usr.Phone, usr.Password)
		log.Println(err.Error())
		return err
	}

	sqlStatement = `UPDATE profile_info SET birthday = $1, gender = $2 WHERE pid = $3;`
	_, err = tx.Exec(sqlStatement, usr.Email, usr.Phone, usr.Password)
	if err != nil {
		log.Println(sqlStatement, usr.Email, usr.Phone, usr.Password)
		log.Println(err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (ur *sqlUserRepository) GetValidTags() ([]models.Tag, error) {
	sqlStatement := `SELECT tag_id, name FROM tag ORDER BY tag_id;`
	rows, err := ur.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	var tags []models.Tag
	defer rows.Close()
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

func (ur *sqlUserRepository) GetRandomFeedUsers(uid int, limit int, page int) ([]models.UserGeneral, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("Page number can't be less than 1\n")
	}
	withCondition := `
		WITH voted_users AS ( SELECT user_id FROM user_vote WHERE uid = $1 ) `
	sqlCondition := ` 
		LEFT JOIN voted_users AS v 
		ON p.pid = v.user_id 
		WHERE v.user_id IS NULL 
		AND
		p.pid != $2
		LIMIT $3 ; `
	// TODO: add cool vote algorithm (aka select)
	return ur.getUsers(withCondition, sqlCondition, uid, uid, limit)
}

// +/- universal method for getting users array by condition (aka sqlStatement)
// and parameters in args (interface array)
func (ur *sqlUserRepository) getUsers(withCondition string, sqlStatement string, args ...interface{}) ([]models.UserGeneral, error) {
	baseSql := withCondition + `
		SELECT p.pid, u.name, p.photos, p.about, p.birthday, p.gender
		FROM profile_info as p
		JOIN profile as u
		ON p.pid = u.uid `
	baseSql += sqlStatement
	log.Println(baseSql, args)
	rows, err := ur.db.Query(baseSql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserGeneral
	for rows.Next() {
		userInfo := models.UserGeneral{}
		gender := ""
		genderPtr := &gender
		bday := time.Time{}
		bdayPtr := &bday
		err = rows.Scan(
			&userInfo.Uid,
			&userInfo.Name,
			&userInfo.Photos,
			&userInfo.About,
			&bdayPtr,
			&genderPtr)
		if err != nil {
			return nil, err
		}
		users = append(users, userInfo)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *sqlUserRepository) GetUserSubscriptions(uid int) ([]models.Event, error) {
	sqlStatement := `
		SELECT
			events.eid, events.uid, title, edate, message, events.is_edited, author, etype, range
		FROM
			events
			JOIN event_vote 
				ON events.eid = event_vote.eid
		WHERE
			event_vote.uid = $1
			AND
			event_vote.value = 1;`
	rows, err := ur.db.Query(sqlStatement, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Event
	for rows.Next() {
		event := models.Event{}
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
		subs = append(subs, event)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return subs, nil
}
