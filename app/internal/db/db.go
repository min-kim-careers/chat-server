package db

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	DBPool *pgxpool.Pool
}

func NewDB(ctx context.Context) *DB {
	db := &DB{
		DBPool: initDBPool(ctx),
	}

	db.RunRoomSchemaSQL(ctx)
	db.RunMessageSchemaSQL(ctx)

	return db
}

func initDBPool(ctx context.Context) *pgxpool.Pool {
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
		// log.Println("Acquiring a DB connection.")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		// log.Println("Releasing a DB connection.")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Database closing.")
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatal("error while creating a DB pool:", err)
	}

	return dbPool
}

//go:embed schema/message.sql
var messageSQL string

func (d *DB) RunMessageSchemaSQL(ctx context.Context) {
	_, err := d.DBPool.Exec(ctx, messageSQL)
	if err != nil {
		log.Fatalf("error creating message table: %v", err)
	}
}

//go:embed schema/room.sql
var roomSQL string

func (d *DB) RunRoomSchemaSQL(ctx context.Context) {
	_, err := d.DBPool.Exec(ctx, roomSQL)
	if err != nil {
		log.Fatalf("error creating room table: %v", err)
	}
}
