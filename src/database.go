package src

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// singleton connection pool
var pool *pgxpool.Pool

func InitConnectionPool() (err error) {

	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)

	if pool, err = pgxpool.New(context.Background(), dsn); err != nil {
		return
	}

	if err = pool.Ping(context.Background()); err != nil {
		return
	}

	return

}

func CloseConnectionPool() {
	pool.Close()
}

func GetConnectionPool() *pgxpool.Pool {
	return pool
}
