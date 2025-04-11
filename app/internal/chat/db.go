package chat

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbCtx = context.Background()

type DB struct {
	pool *pgxpool.Pool
}

func NewDB() *DB {
	return &DB{
		pool: initDBPool(),
	}
}

func initDBPool() *pgxpool.Pool {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing required database environment variables")
	}

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	dbConfig, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		log.Fatal("Failed to create a config:", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Acquiring a DB connection.")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("Releasing a DB connection.")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Database closing.")
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("Error while creating a DB pool:", err)
	}

	return dbPool
}

func (db *DB) AcquireConn() (*pgxpool.Conn, error) {
	dbConn, err := db.pool.Acquire(context.Background())
	if err != nil {
		log.Println("Error acquiring DB connection:", err)
	}
	return dbConn, err
}

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

func (db *DB) Insert(msg *Message) bool {
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

func (db *DB) BulkInsert(msgs []*Message) bool {
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
		parsedTime, err := time.Parse(TIMESTAMP_FORMAT, msg.Timestamp)
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

func (db *DB) Restore(roomID, timestamp string, limit int) []*Message {
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

	msgs := []*Message{}

	for rows.Next() {
		var ts time.Time
		var msg Message
		err := rows.Scan(&msg.Type, &msg.RoomID, &msg.ClientID, &ts, &msg.Content)
		if err != nil {
			log.Printf("Error scanning queried rows from DB for room <%s>: %v", roomID, err)
			return nil
		}
		msg.Timestamp = ts.Format(TIMESTAMP_FORMAT)
		msgs = append(msgs, &msg)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating queried rows from DB for room <%s>: %v", roomID, err)
		return nil
	}

	log.Printf("Fetched %d messages from DB for room <%s>.", len(msgs), roomID)
	return msgs
}
