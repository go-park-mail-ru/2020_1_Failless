package repository

//go:generate mockgen -destination=../mocks/mock_repository.go -package=mocks failless/internal/pkg/chat Repository

import (
	"failless/internal/pkg/chat"
	mydb "failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"log"
)

const (
	QueryInsertGlobalChat = `
		INSERT INTO 	chat_user (admin_id, user_count, title)
		VALUES 			($1, $2, $3)
		RETURNING 		chat_id;`
	QueryInsertFirstMessage = `
		INSERT INTO 	message (uid, chat_id, user_local_id, message, is_shown)
		VALUES 			($1, $2, $3, 'Напишите первое сообщение!', $4);`
	QueryInsertNewLocalChat = `
		INSERT INTO 	user_chat (chat_local_id, uid, avatar, title)
		VALUES			($1, $2, $3, $4)
		RETURNING 		user_local_id;`
	QueryUpdateEidAvatarInChatUser = `
		UPDATE  		chat_user
		SET     		eid = $1,
						avatar = $2
		WHERE   		chat_id = $3;`
	QuerySelectChatInfoByUserId = `
		SELECT 		cu.chat_id, cu.admin_id, cu.date, cu.user_count, cu.title
		FROM 		user_chat uc 
		JOIN 		chat_user cu
		ON 			uc.chat_local_id = cu.chat_id
		WHERE 		uid = $1;`
	QuerySelectMessageHistory = `
		SELECT		ms.mid, ms.uid, ms.chat_id, ms.user_local_id, ms.message, ms.is_shown, ms.created 
		FROM		user_chat uc
		JOIN 		message ms
		ON 			uc.user_local_id = ms.user_local_id
		AND			ms.chat_id = $1
		WHERE		uc.uid = $2 
		ORDER BY	ms.created DESC
		LIMIT		$3
		OFFSET		$4;`
	QuerySelectCheckRoom = `
		SELECT 	cu.chat_id, cu.admin_id, cu.date, cu.user_count, cu.title
		FROM 	user_chat uc 
		JOIN 	chat_user cu ON uc.chat_local_id = $1
		WHERE 	uid = $2;`
	QuerySelectChatsWithLastMsg = `
		SELECT	c.chat_id,
				c.title,
				SUM(CASE WHEN m.is_shown = FALSE THEN 1 ELSE 0 END) AS unseen,
				MAX(m.created) AS last_date,
				SUBSTRING(MAX(m.created || '-----' || m.message) from '%#"-----%#"%' for '#') last_msg,
				uc.avatar,
				uc.title,
				c.user_count
		FROM 	user_chat uc
		JOIN 	chat_user c
		ON 		c.chat_id = uc.chat_local_id
		JOIN 	message m
		ON 		m.user_local_id = uc.user_local_id
		WHERE	uc.uid = $1
		GROUP BY	c.chat_id, uc.avatar, uc.title
		ORDER BY	last_date DESC
		LIMIT	$2
		OFFSET	$3;`
)

type sqlChatRepository struct {
	pgxdb *pgx.ConnPool
	db mydb.MyDBInterface
}

func NewSqlChatRepository(db *pgx.ConnPool) chat.Repository {
	return &sqlChatRepository{
		pgxdb: db,
		db: mydb.NewDBInterface(),
	}
}

func (cr *sqlChatRepository) InsertDialogue(uid1, uid2, userCount int, title string) (int64, error) {
	// Create transaction
	tx, err := cr.pgxdb.Begin()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	defer tx.Rollback()

	// Create global chat
	var chatId int64
	if err := tx.QueryRow(QueryInsertGlobalChat, uid1, userCount, title).Scan(&chatId); err != nil {
		log.Println(err)
		return -1, err
	}

	// Insert into user_chat joining name of person + avatar
	var userLocalIDs [2]int
	sqlStatement := `
		INSERT INTO 	user_chat (chat_local_id, uid, avatar, title)
		SELECT 			$1, $2, pi.photos[1], p.name
		FROM 			profile_info pi
		JOIN 			profile p ON p.uid = pi.pid
		WHERE 			pi.pid = $3
		RETURNING 		user_local_id;`
	if err := tx.QueryRow(sqlStatement, chatId, uid1, uid2).Scan(
		&userLocalIDs[0]); err != nil {
		log.Println(err)
		return -1, err
	}
	if err := tx.QueryRow(sqlStatement, chatId, uid2, uid1).Scan(
		&userLocalIDs[1]); err != nil {
		log.Println(err)
		return -1, err
	}

	// Modify user_vote
	sqlStatement2 := `
		UPDATE 	user_vote
		SET 	chat_id = $1
		WHERE 	(uid = $2 AND user_id = $3)
		OR 		(uid = $3 AND user_id = $2);`
	if row, err := tx.Exec(sqlStatement2, chatId, uid1, uid2); err != nil {
		log.Println(err)
		log.Println(sqlStatement2, chatId, uid1, uid2)
		log.Println(row)
		return -1, err
	}

	//Insert first message
	for _, userLocalID := range userLocalIDs {
		if row, err := tx.Exec(QueryInsertFirstMessage, uid1, chatId, userLocalID, false); err != nil {
			log.Println(err)
			log.Println(row)
			return -1, err
		}
	}

	// Close transaction
	if err = tx.Commit(); err != nil {
		log.Println(err)
		return -1, err
	}

	log.Println("New chat has been created")

	return chatId, nil
}

func (cr *sqlChatRepository) getMessages(sqlStatement string, cid int64, uid int64, limit int, page int) ([]forms.Message, error) {
	rows, err := cr.db.Query(sqlStatement, cid, uid, limit, page)
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
	rows, err := cr.db.Query(QuerySelectChatInfoByUserId, uid)
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
	rows, err := cr.db.Query(QuerySelectCheckRoom, cid, uid)
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
	tx, err := cr.pgxdb.Begin()
	if err != nil {
		return -1, err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()

	sqlStatement := `INSERT INTO message (uid, chat_id, user_local_id, message, is_shown)
						SELECT $1, chat_local_id, user_local_id, $3, $4 FROM 
							(SELECT * FROM user_chat uc WHERE $2 = uc.chat_local_id) AS s1
					RETURNING mid`
	mID := int64(0)
	err = tx.QueryRow(
		sqlStatement,
		msg.Uid,
		msg.ChatID,
		msg.Text,
		true).Scan(&mID)
	if err != nil {
		log.Println("AddMessageToChat: error - ", err.Error())
		log.Println(sqlStatement, msg.Uid, msg.ChatID, msg.ULocalID, true)
		return -1, err
	}
	//sqlStatement = `INSERT INTO message (uid, chat_id, user_local_id, message, is_shown) VALUES `
	//valuesStr := ``
	//values := []interface{}{}
	//itemNum := 5
	//postNum := len(relatedChats)
	//for i := 0; i < postNum; i++ {
	//	valuesStr += ` ( `
	//	for j := 1; j <= itemNum; j++ {
	//		valuesStr += `$` + strconv.Itoa(i*itemNum+j)
	//		if j != itemNum {
	//			valuesStr += `, `
	//		}
	//	}
	//	valuesStr += ` ) `
	//	values = append(values, msg.Uid, msg.ChatID, relatedChats[i], true)
	//	if i != postNum-1 {
	//		valuesStr += ` , `
	//	}
	//}
	//_, err = tx.Exec(sqlStatement+valuesStr+" ;", values...)
	//if err != nil {
	//	return -1, err
	//}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return mID, nil
}

func (cr *sqlChatRepository) GetUserTopMessages(uid int64, page, limit int) ([]models.ChatMeta, error) {
	rows, err := cr.db.Query(QuerySelectChatsWithLastMsg, uid, limit, page)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var chatsMeta []models.ChatMeta
	for rows.Next() {
		meta := models.ChatMeta{}
		ava := pgtype.Varchar{}
		err = rows.Scan(
			&meta.ChatID,
			&meta.Title,
			&meta.Unseen,
			&meta.LastDate,
			&meta.LastMsg,
			&ava,
			&meta.Name,
			&meta.UserCount)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		meta.Avatar = ava.String
		chatsMeta = append(chatsMeta, meta)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return chatsMeta, nil
}

func (cr *sqlChatRepository) GetRoomMessages(uid, cid int64, page, limit int) ([]forms.Message, error) {
	page = 0
	return cr.getMessages(QuerySelectMessageHistory, cid, uid, limit, page)
}
