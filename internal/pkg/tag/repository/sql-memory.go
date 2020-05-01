package repository

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/tag"
	"github.com/jackc/pgx"
)

type sqlTagRepository struct {
	db *pgx.ConnPool
}

func NewSqlTagRepository(db *pgx.ConnPool) tag.Repository {
	return &sqlTagRepository{db: db}
}

func (tr *sqlTagRepository) GetAllTags() ([]models.Tag, error) {
	sqlStatement := `SELECT tag_id, name FROM tag ORDER BY tag_id;`
	rows, err := tr.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		tag := models.Tag{}
		err = rows.Scan(
			&tag.TagId,
			&tag.Name)
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
