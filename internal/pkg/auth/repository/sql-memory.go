package repository

//go:generate mockgen -destination=../mocks/mock_repository.go -package=mocks failless/internal/pkg/auth Repository

import (
	"failless/internal/pkg/auth"
	"failless/internal/pkg/models"
	"github.com/jackc/pgx"
	"log"
)

type sqlAuthRepository struct {
	db *pgx.ConnPool
}

func NewSqlAuthRepository(db *pgx.ConnPool) auth.Repository {
	return &sqlAuthRepository{db: db}
}

// Private method
func (ur *sqlAuthRepository) getUser(sqlStatement string, args ...interface{}) (user models.User, err error) {
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

func (ur *sqlAuthRepository) GetUserByPhoneOrEmail(phone string, email string) (models.User, error) {
	sqlStatement := `SELECT uid, name, phone, email, password FROM profile WHERE phone = $1 OR email = LOWER($2);`
	return ur.getUser(sqlStatement, phone, email)
}
