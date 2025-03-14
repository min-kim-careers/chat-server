package db

import (
	"chat-go/main/chat"
	"context"
	"fmt"
	"log"
	"os"

	"chat-app/main/chat"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getDBConnStr() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing required database environment variables")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func createMessagesTable(dbpool *pgxpool.Conn) {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender VARCHAR(255) NOT NULL,
		receiver VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = dbpool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to apply schema: %v", err)
		return nil, err
	}
}

func InitDB() (*pgxpool.Conn, error) {
	dbConnStr := getDBConnStr()

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	

	return dbpool, nil
}

func AddMessage(dbConn *pgx.Conn, msg *chat.Message) {
	dbConn.
}