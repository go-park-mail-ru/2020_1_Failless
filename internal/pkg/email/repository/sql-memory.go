package repository

import (
	"failless/internal/pkg/email"
	"failless/internal/pkg/models"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

const (
	QueryInsertEmail = `
		INSERT INTO	emails (email)
		VALUES		($1);`
)

type sqlEmailsRepository struct {
	db *pgx.ConnPool
}

func NewSqlEmailRepository(db *pgx.ConnPool) email.Repository {
	return &sqlEmailsRepository{db: db}
}

func (er *sqlEmailsRepository) SaveEmail(email *models.Email) (int, error) {
	cTag, err := er.db.Exec(QueryInsertEmail, email.Email)
	if err != nil || cTag.RowsAffected() == 0 {
		log.Println(err)
		return http.StatusInternalServerError, err
	} else {
		return http.StatusOK, nil
	}
}
