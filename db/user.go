package db

import (
	"errors"
	"github.com/jackc/pgx"
	"log"
)

func GetUserByPhoneOrEmail(db *pgx.ConnPool, phone string, email string) (User, error) {
	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE phone = $1 OR email = $2;`
	row := db.QueryRow(sqlStatement, phone, email)

	var user User
	err := row.Scan(
		&user.Uid,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Password)
	if err == pgx.ErrNoRows {
		return User{-1, "", "", "", []byte{}}, nil
	} else if err != nil {
		return User{}, err
	}
	return user, nil
}

func AddNewUser(db *pgx.ConnPool, user *User) error {
	sqlStatement := `INSERT INTO profile VALUES (default, $1, $2, $3, $4) RETURNING uid;`
	uid := int(0)
	err := db.QueryRow(sqlStatement, user.Name, user.Phone, user.Email, user.Password).Scan(&uid)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	user.Uid = uid
	sqlStatement = `INSERT INTO profile_info VALUES ( $1 , default , default , default , default , default , default , default ) ;`
	_, err = db.Exec(sqlStatement, user.Uid)
	return err
}

func AddUserInfo(db *pgx.ConnPool, user User, info UserInfo) error {
	sqlStatement := `SELECT uid FROM profile WHERE LOWER(email) = LOWER($1) OR phone = $2;`
	row := db.QueryRow(sqlStatement, user.Email, user.Phone)
	err := row.Scan(&user.Uid)
	if err == pgx.ErrNoRows {
		return errors.New("user " + user.Email + "doesn't exist")
	} else if err != nil {
		return err
	}
	sqlStatement = `UPDATE profile_info SET about = $1, photos = $2, birthday = $3, gender = $4 WHERE pid = $5;`
	_, err = db.Exec(sqlStatement, info.About, info.Photos, info.Birthday, info.Gender, user.Uid)
	if err != nil {
		return err
	}
	return nil
}

func SetUserLocation(db *pgx.ConnPool, uid int, point LocationPoint) error {
	sqlStatement := `UPDATE profile_info SET location = ST_POINT($1, $2) WHERE pid = $3;`
	_, err := db.Exec(sqlStatement, point.Latitude, point.Longitude, uid)
	return err
}

func UpdateUserRating(db *pgx.ConnPool, uid int, rating float32) error {
	sqlStatement := `UPDATE profile_info SET rating = $1 WHERE pid = $2;`
	_, err := db.Exec(sqlStatement, rating, uid)
	return err
}

func GetProfileInfo(db *pgx.ConnPool, uid int) (user UserInfo, err error) {
	sqlStatement := `SELECT about, photos, rating, birthday, gender FROM profile_info WHERE pid = $1 ;`
	err = db.QueryRow(sqlStatement, uid).Scan(
		&user.About,
		&user.Photos,
		&user.Rating,
		&user.Birthday,
		&user.Gender)
	if err != nil {
		log.Println(sqlStatement)
		log.Println("error in get profile info")
		return UserInfo{}, err
	}
	return
}
