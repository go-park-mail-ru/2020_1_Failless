package db

import "github.com/jackc/pgx"

func GetUserByPhoneOrEmail(db *pgx.ConnPool, phone string, email string) (UserInfo, error) {
	var user User
	var userInfo UserInfo

	sqlStatement := `SELECT uid, full_name, nickname, email, about FROM profile WHERE nickname = $1;`
	row := db.QueryRow(sqlStatement, phone, email)
	err := row.Scan(
		&userInfo.Pk,
		&userInfo.Name,
		&userInfo.Nickname,
		&userInfo.Email,
		&userInfo.About)
	if err == pgx.ErrNoRows {
		return UserInfo{}, err
	}
	return userInfo, nil
}
