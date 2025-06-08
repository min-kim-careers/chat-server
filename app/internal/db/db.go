package db

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
