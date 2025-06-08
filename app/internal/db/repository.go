package db

import (
	"chat-server/internal/models"
	"chat-server/internal/utils"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func (db *DB) CreateMessageTable() error {
	dbConn, _ := db.AcquireConn()
	defer dbConn.Release()

	q := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		message_type VARCHAR(255) NOT NULL,
		room_id VARCHAR(255) NOT NULL,
		client_id VARCHAR(255) NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		message_content TEXT NOT NULL
	)
	`

	_, err := dbConn.Exec(context.Background(), q)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	log.Println("Message table found.")
	return nil
}

func (db *DB) Insert(msg *models.Message) bool {
	dbConn, err := db.AcquireConn()
	defer dbConn.Release()
	if err != nil {
		return false
	}

	q := `
	INSERT INTO messages (message_type, room_id, client_id, timestamp, message_content) 
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err = dbConn.Exec(dbCtx, q, msg.Type, msg.RoomID, msg.ClientID, msg.Timestamp, msg.Content)
	if err != nil {
		log.Println("Error adding message:", err)
		return false
	}

	log.Println("Message added to DB.")
	return true
}

func (db *DB) BulkInsert(msgs []*models.Message) bool {
	if len(msgs) == 0 {
		log.Println("No messages to persist")
		return false
	}

	dbConn, err := db.AcquireConn()
	defer dbConn.Release()
	if err != nil {
		return false
	}

	rows := make([][]any, len(msgs))
	for i, msg := range msgs {
		parsedTime, err := time.Parse(utils.TIMESTAMP_FORMAT, msg.Timestamp)
		if err != nil {
			log.Println("Error parsing timestamp:", err)
			return false
		}

		rows[i] = []any{msg.Type, msg.RoomID, msg.ClientID, parsedTime, msg.Content}
	}

	_, err = dbConn.CopyFrom(
		dbCtx,
		pgx.Identifier{"messages"},
		[]string{"message_type", "room_id", "client_id", "timestamp", "message_content"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Printf("Error bulk inserting messages: %v", err)
		return false
	}

	log.Println("Successfully bulk inserted messages")
	return true
}

func (db *DB) Restore(roomID, timestamp string, limit int) []*models.Message {
	if roomID == "" || timestamp == "" || limit <= 0 {
		log.Println("Invalid parameters passed to Restore")
		return nil
	}

	dbConn, err := db.AcquireConn()
	defer dbConn.Release()
	if err != nil {
		log.Printf("Error restoring DB messages for room <%s>.", roomID)
		return nil
	}

	q := `
	SELECT message_type, room_id, client_id, timestamp, message_content FROM messages 
	WHERE room_id = $1 AND timestamp < $2
	ORDER BY timestamp DESC
	LIMIT $3
	`

	rows, err := dbConn.Query(dbCtx, q, roomID, timestamp, limit)
	if err != nil {
		log.Printf("Error querying messages from DB for room <%s>: %v", roomID, err)
		return nil
	}
	defer rows.Close()

	msgs := []*models.Message{}

	for rows.Next() {
		var ts time.Time
		var msg models.Message
		err := rows.Scan(&msg.Type, &msg.RoomID, &msg.ClientID, &ts, &msg.Content)
		if err != nil {
			log.Printf("Error scanning queried rows from DB for room <%s>: %v", roomID, err)
			return nil
		}
		msg.Timestamp = ts.Format(utils.TIMESTAMP_FORMAT)
		msgs = append(msgs, &msg)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating queried rows from DB for room <%s>: %v", roomID, err)
		return nil
	}

	log.Printf("Fetched %d messages from DB for room <%s>.", len(msgs), roomID)
	return msgs
}
