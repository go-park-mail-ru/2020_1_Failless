package db

import "github.com/jackc/pgx"

func GetUserByPhoneOrEmail(db *pgx.ConnPool, phone string, email string) (User, error) {
	var user User

	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE phone = $1 OR email = $2;`
	row := db.QueryRow(sqlStatement, phone, email)
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
