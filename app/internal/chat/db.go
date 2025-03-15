package chat

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbCtx = context.Background()

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() *DB {
	db := DB{}
	dbPool, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Fatal("Error while creating pool to the database:", err)
	}
	db.Pool = dbPool
	return &db
}

func (db *DB) AcquireConn() (*pgxpool.Conn, error) {
	dbConn, err := db.Pool.Acquire(context.Background())
	if err != nil {
		log.Println("Error acquiring connection from pool:", err)
	}
	return dbConn, err
}

func (db *DB) CreateMessageTable() {
	dbConn, _ := db.AcquireConn()
	defer dbConn.Release()

	q := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		message_type VARCHAR(255) NOT NULL,
		room_id VARCHAR(255) NOT NULL,
		client_id VARCHAR(255) NOT NULL,
		message_content TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := dbConn.Exec(context.Background(), q)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

func (db *DB) AddMessage(msg *Message) {
	dbConn, err := db.AcquireConn()
	defer dbConn.Release()
	if err != nil {
		return
	}

	q := `INSERT INTO messages (message_type, room_id, client_id, message_content, timestamp) VALUES ($1, $2, $3, $4, $5)`

	_, err = dbConn.Exec(dbCtx, q, msg.Type, msg.RoomID, msg.ClientID, msg.Content, msg.Timestamp)
	if err != nil {
		log.Println("Error adding message:", err)
	} else {
		log.Println("Message added to DB.")
	}
}

func (db *DB) GetMessageHistory(roomID RoomID) {
	dbConn, err := db.AcquireConn()
	defer dbConn.Release()
	if err != nil {
		return
	}

	// q := `SELECT * FROM messages WHERE room_id = $1`

	// rows, err := dbConn.Query(dbCtx, q, roomID)
	// if err != nil {
	// 	log.Printf("Error fetching message history for room <%s>.", roomID)
	// } else {
	// 	log.Printf("Fetched message history for room <%s>.", roomID)
	// }
}
