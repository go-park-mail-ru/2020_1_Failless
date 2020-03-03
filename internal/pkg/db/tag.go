package db

import "github.com/jackc/pgx"

func GetAllTags(db *pgx.ConnPool) ([]Tag, error) {
	sqlStatement := `SELECT tag_id, name FROM tag ORDER BY tag_id;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	var tags []Tag
	for rows.Next() {
		tag := Tag{}
		err = rows.Scan(
			&tag.Name,
			&tag.TagId)
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
