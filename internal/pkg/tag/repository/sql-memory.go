package repository

//go:generate mockgen -destination=../mocks/mock_repository.go -package=mocks failless/internal/pkg/tag Repository

import (
	mydb "failless/internal/pkg/db"
	"failless/internal/pkg/models"
	"failless/internal/pkg/tag"
)

const (
	QuerySelectTagById = `
		SELECT 		tag_id, name
		FROM 		tag
		ORDER BY 	tag_id;`
)

type sqlTagRepository struct {
	db mydb.MyDBInterface
}

func NewSqlTagRepository() tag.Repository {
	return &sqlTagRepository{db: mydb.NewDBInterface()}
}

func (tr *sqlTagRepository) GetAllTags() ([]models.Tag, error) {
	rows, err := tr.db.Query(QuerySelectTagById)
	if err != nil {
		return nil, err
	}

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
