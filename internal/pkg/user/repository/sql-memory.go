package repository

//go:generate mockgen -destination=../mocks/mock_repository.go -package=mocks failless/internal/pkg/user Repository

import (
	"errors"
	mydb "failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/user"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"log"
	"net/http"
	"time"
)

const (
	QueryUpdateUserAbout = `
		UPDATE  profile_info
		SET     about = $1
		WHERE   pid = $2;`
	QueryUpdateUserTags = `
		UPDATE  profile_info
		SET     tags = $1
		WHERE   pid = $2;`
	QueryUpdateUserPhotos = `
		UPDATE 	profile_info
		SET 	photos = $1
		WHERE 	pid = $2;`
	QuerySelectUserByID = `
		SELECT 	uid, name, phone, email, password
		FROM 	profile
		WHERE 	uid = $1;`
	QuerySelectUserByPhoneOrEmail = `
		SELECT 	uid, name, phone, email, password
		FROM 	profile
		WHERE 	phone = $1
		OR 		email = LOWER($2);`
	QueryUpdateUserLocation = `
		UPDATE 	profile_info
		SET 	location = ST_POINT($1, $2)
		WHERE 	pid = $3;`
	QueryUpdateUserRating = `
		UPDATE 	profile_info
		SET 	rating = $1
		WHERE 	pid = $2;`
	QueryDeleteUserByEmail = `
		DELETE FROM 	profile
		WHERE 			email=$1;`
	QuerySelectTags = `
		SELECT 		tag_id, name
		FROM 		tag
		ORDER BY 	tag_id;`
	QuerySelectUserInfoIncomplete = `
		SELECT	p.pid, u.name, p.photos, p.about, p.birthday, p.gender, p.tags
		FROM	profile_info as p
		JOIN	profile as u
		ON		p.pid = u.uid `
	QueryWithVotedUsersIncomplete = `
		WITH	voted_users AS (SELECT user_id FROM user_vote WHERE uid = $1 ) `
	QueryConditionFeedIncomplete = ` 
		LEFT JOIN	voted_users AS v
		ON			p.pid = v.user_id
		WHERE		v.user_id IS NULL
		AND			p.pid != $1
		AND			p.photos IS NOT NULL	-- we don't show users without full profile info
		AND			p.about IS NOT NULL
		AND			p.about <> ''
		LIMIT		$2;`
	QueryWithChatMembersIncomplete = `
		WITH members_id AS (
			SELECT  mem.uid
			FROM    mid_events
			JOIN    mid_event_members mem
			ON      mid_events.eid = mem.eid
			WHERE   chat_id = $1
		)`
	QueryConditionChatMembersIncomplete = `
		JOIN		members_id AS mi
		ON			p.pid = mi.uid;`
)

var (
	CorrectMessage = models.WorkMessage{
		Request: nil,
		Message: "",
		Status:  http.StatusOK,
	}
)

type sqlUserRepository struct {
	pgxdb *pgx.ConnPool
	db mydb.MyDBInterface
}

func NewSqlUserRepository(db *pgx.ConnPool) user.Repository {
	return &sqlUserRepository{
		pgxdb: db,
		db: mydb.NewDBInterface(),
	}
}

func (ur *sqlUserRepository) GetUserByUID(uid int) (models.User, error) {
	return ur.getUser(QuerySelectUserByID, uid)
}

func (ur *sqlUserRepository) GetUserByPhoneOrEmail(phone string, email string) (models.User, error) {
	return ur.getUser(QuerySelectUserByPhoneOrEmail, phone, email)
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

func (ur *sqlUserRepository) AddNewUser(user *models.User) error {
	uid := 0
	sqlStatement := `
		INSERT INTO profile
		VALUES 		(default, $1, $2, LOWER($3), $4)
		RETURNING 	uid;`
	err := ur.db.QueryRow(sqlStatement, user.Name, user.Phone, user.Email, user.Password).Scan(&uid)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	user.Uid = uid
	sqlStatement = `
		INSERT INTO profile_info
		VALUES ( $1 , '' , default , default , default , default , default , default ) ;`
	_, err = ur.db.Exec(sqlStatement, user.Uid)
	if err != nil {
		log.Println(sqlStatement, user.Uid)
		return err
	}
	return nil
}

func (ur *sqlUserRepository) SetUserLocation(uid int, point models.LocationPoint) error {
	_, err := ur.db.Exec(QueryUpdateUserLocation, point.Latitude, point.Longitude, uid)
	return err
}

func (ur *sqlUserRepository) UpdateUserRating(uid int, rating float32) error {
	_, err := ur.db.Exec(QueryUpdateUserRating, rating, uid)
	return err
}

func (ur *sqlUserRepository) UpdateUserTags(uid int, tagIDs []int) models.WorkMessage {
	cTag, err := ur.db.Exec(QueryUpdateUserTags, tagIDs, uid)
	if err != nil || cTag.RowsAffected() == 0 {
		log.Println(err)
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	} else {
		return CorrectMessage
	}
}

func (ur *sqlUserRepository) UpdateUserAbout(uid int, about string) models.WorkMessage {
	cTag, err := ur.db.Exec(QueryUpdateUserAbout, about, uid)
	if err != nil || cTag.RowsAffected() == 0 {
		log.Println(err)
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	} else {
		return CorrectMessage
	}
}

func (ur *sqlUserRepository) GetProfileInfo(uid int) (info models.JsonInfo, err error) {
	var profile ProfileInfo
	sqlStatement := `
		SELECT 	about, photos, rating, birthday, gender, tags
		FROM 	profile_info
		WHERE 	pid = $1 ;`
	gender := ""
	genPtr := &gender
	row := ur.db.QueryRow(sqlStatement, uid)
	tags := pgtype.Int4Array{}
	err = row.Scan(
		&profile.About,
		&info.Photos,
		&profile.Rating,
		&profile.Birthday,
		&genPtr,
		&tags)
	if err != nil {
		log.Println(sqlStatement)
		log.Println("error in get profile info")
		return models.JsonInfo{}, err
	}

	if profile.About != nil {
		info.About = *profile.About
	}

	if genPtr != nil {
		info.Gender = user.GenderByStr(gender)
	}

	if profile.Location != nil {
		info.Location = *profile.Location
	}

	if err = tags.AssignTo(&info.Tags); err != nil {
		log.Println("GetProfileInfo: Assigning tags:", err)
	}
	return
}

// func for testing
func (ur *sqlUserRepository) DeleteUser(mail string) error {
	_, err := ur.db.Exec(QueryDeleteUserByEmail, mail)
	return err
}

func (ur *sqlUserRepository) UpdateUserPhotos(uid int, newImages *[]string) models.WorkMessage {
	if cTag, err := ur.db.Exec(QueryUpdateUserPhotos, &newImages, uid); err != nil || cTag.RowsAffected() == 0 {
		log.Println(err)
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return CorrectMessage
}

func (ur *sqlUserRepository) UpdUserGeneral(info models.JsonInfo, usr models.User) error {
	//gender := user.GenderById(info.Gender)

	tx, err := ur.pgxdb.Begin()
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()

	// TODO: add name
	sqlStatement := `
		UPDATE 	profile
		SET 	email = $1, phone = $2, password = $3
		WHERE 	uid = $4;`
	_, err = tx.Exec(sqlStatement, usr.Email, usr.Phone, usr.Password)
	if err != nil {
		log.Println(sqlStatement, usr.Email, usr.Phone, usr.Password)
		log.Println(err.Error())
		return err
	}

	sqlStatement = `
		UPDATE 	profile_info
		SET 	birthday = $1, gender = $2
		WHERE 	pid = $3;`
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
	rows, err := ur.db.Query(QuerySelectTags)
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

func (ur *sqlUserRepository) GetRandomFeedUsers(uid int, limit int, page int) ([]models.UserGeneral, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("Page number can't be less than 1\n")
	}
	return ur.getUsers(QueryWithVotedUsersIncomplete, QueryConditionFeedIncomplete, uid, limit)
}

// +/- universal method for getting users array by condition (aka sqlStatement)
// and parameters in args (interface array)
func (ur *sqlUserRepository) getUsers(withCondition string, sqlStatement string, args ...interface{}) ([]models.UserGeneral, error) {
	baseSql := withCondition + QuerySelectUserInfoIncomplete + sqlStatement
	rows, err := ur.db.Query(baseSql, args...)
	if err != nil {
		log.Println("getUsers: ", err)
		log.Println(rows)
		return nil, err
	}

	var users []models.UserGeneral
	for rows.Next() {
		tags := pgtype.Int4Array{}
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
			&genderPtr,
			&tags)
		if err != nil {
			return nil, err
		}
		if err = tags.AssignTo(&userInfo.TagsId); err != nil {
			log.Println("EventRepo: GetMidEventsWithFollowed: Assigning tags:", err)
		}
		users = append(users, userInfo)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *sqlUserRepository) GetUsersForChat(cid int64, users *models.UserGeneralList) models.WorkMessage {
	var err error
	*users, err = ur.getUsers(QueryWithChatMembersIncomplete, QueryConditionChatMembersIncomplete, cid)
	if err != nil {
		return models.WorkMessage{
			Request: nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	} else {
		return models.WorkMessage{
			Request: nil,
			Message: "",
			Status:  http.StatusOK,
		}
	}
}
