package repository

import (
	"failless/internal/pkg/chat"
	"github.com/jackc/pgx"
	"log"
)

type sqlChatRepository struct {
	db *pgx.ConnPool
}

func NewSqlChatRepository(db *pgx.ConnPool) chat.Repository {
	return &sqlChatRepository{db: db}
}

func (cr *sqlChatRepository) InsertDialogue(id1, id2 int) (int, error) {
	// Create transaction
	tx, err := cr.db.Begin()
	if err != nil {
		log.Println(err)
		return -1, nil
	}
	defer tx.Rollback()

	// Add row to chat_pair
	var chatId int
	sqlStatement1 := `
		INSERT INTO chat_pair (id1, id2)
		VALUES ($1, $2)
		RETURNING chat_id;`
	if err := tx.QueryRow(sqlStatement1, id1, id2).Scan(&chatId); err != nil {
		log.Println(err)
		return -1, nil
	}

	// Modify user_vote
	sqlStatement2 := `
		UPDATE user_vote
		SET chat_id = $1
		WHERE (uid = $2 AND user_id = $3) OR (uid = $3 AND user_id = $2);`
	if _, err = tx.Exec(sqlStatement2, chatId, id1, id2); err != nil {
		log.Println(err)
		return -1, nil
	}

	// Close transaction
	if err = tx.Commit(); err != nil {
		log.Println(err)
		return -1, err
	}

	return chatId, nil
}
