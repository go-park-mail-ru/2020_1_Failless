package db

import (
	"errors"
	"github.com/jackc/pgx"
	"log"
)

func GetUserByPhoneOrEmail(db *pgx.ConnPool, phone string, email string) (User, error) {
	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE phone = $1 OR email = LOWER($2);`
	return getUser(db, sqlStatement, phone, email)
}

func GetUserByUID(db *pgx.ConnPool, uid int) (User, error) {
	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE uid = $1;`
	return getUser(db, sqlStatement, uid)
}

func getUser(db *pgx.ConnPool, sqlStatement string, args... interface{}) (user User, err error) {
	row := db.QueryRow(sqlStatement, args...)
	err = row.Scan(
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
	log.Println(user)
	return user, nil
}



func AddNewUser(db *pgx.ConnPool, user *User) error {
	uid := 0
	sqlStatement := `INSERT INTO profile VALUES (default, $1, $2, LOWER($3), $4) RETURNING uid;`
	err := db.QueryRow(sqlStatement, user.Name, user.Phone, user.Email, user.Password).Scan(&uid)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(sqlStatement, user.Name, uid)
	user.Uid = uid
	sqlStatement = `INSERT INTO profile_info VALUES ( $1 , 'Расскажите о себе' , default , default , default , default , default , default ) ;`
	_, err = db.Exec(sqlStatement, user.Uid)
	if err != nil {
		log.Println(sqlStatement, user.Uid)
		return err
	}
	return nil
}

func AddUserInfo(db *pgx.ConnPool, user User, info UserInfo) error {
	if user.Uid < 0 {
		sqlStatement := `SELECT uid FROM profile WHERE LOWER(email) = LOWER($1) OR phone = $2;`
		row := db.QueryRow(sqlStatement, user.Email, user.Phone)
		err := row.Scan(&user.Uid)
		if err == pgx.ErrNoRows {
			return errors.New("user " + user.Email + "doesn't exist")
		} else if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	gender := "other"
	switch info.Gender {
	case 0:
		gender = "male"
	case 1:
		gender = "female"
	}

	sqlStatement := `UPDATE profile_info SET about = $1, photos = array_append(photos, $2), birthday = $3, gender = $4 WHERE pid = $5;`
	_, err := db.Exec(sqlStatement, info.About, info.Photos[0], info.Birthday, gender, user.Uid)
	if err != nil {
		log.Println(sqlStatement, info.About, info.Photos[0], info.Birthday, info.Gender, user.Uid)
		log.Println(err.Error())
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
	var profile ProfileInfo
	sqlStatement := `SELECT about, photos, rating, birthday, gender FROM profile_info WHERE pid = $1 ;`
	gender := ""
	genPtr := &gender
	err = db.QueryRow(sqlStatement, uid).Scan(
		&profile.About,
		&profile.Photos,
		&profile.Rating,
		&profile.Birthday,
		&genPtr)
	if err != nil {
		log.Println(sqlStatement)
		log.Println("error in get profile info")
		return UserInfo{}, err
	}
	if profile.About != nil {
		user.About = *profile.About
	}
	if profile.Photos != nil {
		user.Photos = *profile.Photos
	}
	if genPtr != nil {
		switch gender {
		case "male":
			user.Gender = 0
		case "female":
			user.Gender = 1
		default:
			user.Gender = 2
		}
	}
	if profile.Location != nil {
		user.Location = *profile.Location
	}
	return
}

// func for testing
func DeleteUser(db *pgx.ConnPool, mail string) error {
	sqlStatement := `DELETE FROM profile WHERE email=$1;`
	_, err := db.Exec(sqlStatement, mail)
	return err
}