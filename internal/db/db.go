package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

const defaultDatabaseURL = "postgres://day_planner:day_planner@localhost:5432/day_planner?sslmode=disable"

func Open(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		databaseURL = os.Getenv("DAY_PLANNER_DATABASE_URL")
	}
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}
