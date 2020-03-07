package repository

import (
	"errors"
	"failless/internal/pkg/models"
	"failless/internal/pkg/user"
	"github.com/jackc/pgx"
	"log"
)

type sqlUserRepository struct {
	db *pgx.ConnPool
}

func NewSqlUserRepository(db *pgx.ConnPool) user.Repository {
	return &sqlUserRepository{db: db}
}

func (um *sqlUserRepository) GetUserByUID(uid int) (models.User, error) {
	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE uid = $1;`
	return um.getUser(sqlStatement, uid)
}

func (um *sqlUserRepository) GetUserByPhoneOrEmail(phone string, email string) (models.User, error) {
	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE phone = $1 OR email = LOWER($2);`
	return um.getUser(sqlStatement, phone, email)
}

func (um *sqlUserRepository) getUser(sqlStatement string, args ...interface{}) (user models.User, err error) {
	row := um.db.QueryRow(sqlStatement, args...)
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

func (um *sqlUserRepository) AddNewUser(user *models.User) error {
	uid := 0
	sqlStatement := `INSERT INTO profile VALUES (default, $1, $2, LOWER($3), $4) RETURNING uid;`
	err := um.db.QueryRow(sqlStatement, user.Name, user.Phone, user.Email, user.Password).Scan(&uid)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(sqlStatement, user.Name, uid)
	user.Uid = uid
	sqlStatement = `INSERT INTO profile_info VALUES ( $1 , 'Расскажите о себе' , default , default , default , default , default , default ) ;`
	_, err = um.db.Exec(sqlStatement, user.Uid)
	if err != nil {
		log.Println(sqlStatement, user.Uid)
		return err
	}
	return nil
}

func (um *sqlUserRepository) AddUserInfo(credentials models.User, info models.JsonInfo) error {
	if credentials.Uid < 0 {
		sqlStatement := `SELECT uid FROM profile WHERE LOWER(email) = LOWER($1) OR phone = $2;`
		row := um.db.QueryRow(sqlStatement, credentials.Email, credentials.Phone)
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
	_, err := um.db.Exec(sqlStatement, info.About, info.Photos[0], info.Birthday, gender, credentials.Uid)
	if err != nil {
		log.Println(sqlStatement, info.About, info.Photos[0], info.Birthday, info.Gender, credentials.Uid)
		log.Println(err.Error())
		return err
	}
	return nil
}

func (um *sqlUserRepository) SetUserLocation(uid int, point models.LocationPoint) error {
	sqlStatement := `UPDATE profile_info SET location = ST_POINT($1, $2) WHERE pid = $3;`
	_, err := um.db.Exec(sqlStatement, point.Latitude, point.Longitude, uid)
	return err
}

func (um *sqlUserRepository) UpdateUserRating(uid int, rating float32) error {
	sqlStatement := `UPDATE profile_info SET rating = $1 WHERE pid = $2;`
	_, err := um.db.Exec(sqlStatement, rating, uid)
	return err
}

func (um *sqlUserRepository) GetProfileInfo(uid int) (info models.JsonInfo, err error) {
	var profile ProfileInfo
	sqlStatement := `SELECT about, photos, rating, birthday, gender FROM profile_info WHERE pid = $1 ;`
	gender := ""
	genPtr := &gender
	err = um.db.QueryRow(sqlStatement, uid).Scan(
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
func (um *sqlUserRepository) DeleteUser(mail string) error {
	sqlStatement := `DELETE FROM profile WHERE email=$1;`
	_, err := um.db.Exec(sqlStatement, mail)
	return err
}
