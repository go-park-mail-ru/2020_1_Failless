package repository

import (
	"failless/internal/pkg/vote"
	"github.com/jackc/pgx"
)

type sqlVoteRepository struct {
	db *pgx.ConnPool
}

func NewSqlVoteRepository(db *pgx.ConnPool) vote.Repository {
	return &sqlVoteRepository{db: db}
}
