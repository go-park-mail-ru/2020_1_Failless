package repository

import (
	"failless/internal/pkg/email"
	"github.com/jackc/pgx"
)

type sqlEmailsRepository struct {
	db *pgx.ConnPool
}

func NewSqlEmailRepository(db *pgx.ConnPool) email.Repository {
	return &sqlEmailsRepository{db: db}
}
