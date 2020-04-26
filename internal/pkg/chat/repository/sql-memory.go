package repository

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"github.com/jackc/pgx"
	"log"
	"strconv"
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

func (cr *sqlChatRepository) getMessages(sqlStatement string, args ...interface{}) ([]forms.Message, error) {
	rows, err := cr.db.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	var messages []forms.Message
	for rows.Next() {
		msg := forms.Message{}
		err = rows.Scan(
			&msg.Mid,
			&msg.Uid,
			&msg.ChatID,
			&msg.ULocalID,
			&msg.Text,
			&msg.IsShown,
			&msg.Date)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (cr *sqlChatRepository) GetUsersRooms(uid int64) ([]models.ChatRoom, error) {
	sqlStatement := `SELECT cu.chat_id, cu.admin_id, cu.date, cu.user_count, cu.title FROM user_chat uc 
						JOIN chat_user cu ON uc.chat_local_id = cu.chat_id WHERE uid = $1;`
	rows, err := cr.db.Query(sqlStatement, uid)
	if err != nil {
		return nil, err
	}
	var rooms []models.ChatRoom
	for rows.Next() {
		room := models.ChatRoom{}
		err = rows.Scan(
			&room.ChatID,
			&room.AdminID,
			&room.Created,
			&room.UsersCount,
			&room.Title)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (cr *sqlChatRepository) CheckRoom(cid int64, uid int64) (bool, error) {
	sqlStatement := `SELECT cu.chat_id, cu.admin_id, cu.date, cu.user_count, cu.title FROM user_chat uc 
						JOIN chat_user cu ON uc.chat_local_id = $1 WHERE uid = $2;`
	rows, err := cr.db.Query(sqlStatement, cid, uid)
	if err != nil && rows != nil && !rows.Next() {
		log.Println("CheckRoom: user has no rooms")
		return false, nil
	} else if err != nil {
		log.Println("CheckRoom: error - ", err.Error())
		return false, err
	}
	return true, nil
}

func (cr *sqlChatRepository) AddMessageToChat(msg *forms.Message, relatedChats []int64) (int64, error) {
	tx, err := cr.db.Begin()
	if err != nil {
		return -1, err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()

	sqlStatement := `INSERT INTO message (uid, chat_id, user_local_id, message, is_shown) 
							VALUES ($1, $2, $3, $4, $5) RETURNS mid;`
	mID := int64(0)
	err = tx.QueryRow(
		sqlStatement,
		msg.Uid,
		msg.ChatID,
		msg.ULocalID,
		msg.Text,
		true).Scan(&mID)
	if err != nil {
		log.Println("AddMessageToChat: error - ", err.Error())
		log.Println(sqlStatement, msg.Uid, msg.ChatID, msg.ULocalID, true)
		return -1, err
	}
	sqlStatement = `INSERT INTO message (uid, chat_id, user_local_id, message, is_shown) VALUES `
	valuesStr := ``
	values := []interface{}{}
	itemNum := 5
	postNum := len(relatedChats)
	for i := 0; i < postNum; i++ {
		valuesStr += ` ( `
		for j := 1; j <= itemNum; j++ {
			valuesStr += `$` + strconv.Itoa(i*itemNum+j)
			if j != itemNum {
				valuesStr += `, `
			}
		}
		valuesStr += ` ) `
		values = append(values, msg.Uid, msg.ChatID, relatedChats[i], true)
		if i != postNum-1 {
			valuesStr += ` , `
		}
	}
	_, err = tx.Exec(sqlStatement+valuesStr+" ;", values...)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return mID, nil
}

func (cr *sqlChatRepository) GetUserTopMessages(uid int64, page, limit int) ([]models.ChatMeta, error) {
	sqlStatement := `SELECT c.chat_id,, c.title, SUM(CASE WHEN m.is_shown = TRUE THEN 1 ELSE 0 END) AS unseen,
						MAX(m.created) AS last_date, SUBSTR(MAX(CONCAT(m.created, m.message)), 20) last_msg
						FROM user_chat uc JOIN chat_user c ON c.chat_id = uc.chat_local_id 
						JOIN message m ON m.user_local_id = uc.user_local_id WHERE uc.uid = $1
							GROUP BY c.chat_id ORDER BY last_date ASC LIMIT $2 OFFSET $3;`
	rows, err := cr.db.Query(sqlStatement, uid, limit, page)
	if err != nil {
		return nil, err
	}
	var chatsMeta []models.ChatMeta
	for rows.Next() {
		meta := models.ChatMeta{}
		err = rows.Scan(
			&meta.ChatID,
			&meta.Title,
			&meta.Unseen,
			&meta.LastDate,
			&meta.LastMsg)
		if err != nil {
			return nil, err
		}
		chatsMeta = append(chatsMeta, meta)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return chatsMeta, nil
}

func (cr *sqlChatRepository) GetRoomMessages(uid, cid int64, page, limit int) ([]forms.Message, error) {
	sqlStatement := `SELECT ms.mid, ms.uid, ms.chat_id, ms.user_local_id, ms.message, ms.is_shown, ms.created 
						FROM user_chat uc JOIN message ms ON uc.user_local_id = ms.user_local_id AND ms.chat_id = $1
						WHERE uc.uid = $2 ORDER BY ms.created ASC LIMIT $3 OFFSET $4;`
	return cr.getMessages(sqlStatement, cid, uid, limit, page)
}
