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

func (db *DB) CreateMessageTable() {
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
}

func (db *DB) AddMessage(msg *Message) error {
	dbConn, err := db.AcquireConn()
	defer dbConn.Release()
	if err != nil {
		return err
	}

	q := `
	INSERT INTO messages (message_type, room_id, client_id, timestamp, message_content) 
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err = dbConn.Exec(dbCtx, q, msg.Type, msg.RoomID, msg.ClientID, msg.Timestamp, msg.Content)
	if err != nil {
		log.Println("Error adding message:", err)
	} else {
		log.Println("Message added to DB.")
	}
	return err

}

func (db *DB) AddMessages(msgs []*Message) {

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
