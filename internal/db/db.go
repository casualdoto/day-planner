package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

var ErrDatabaseConfigRequired = errors.New("database config is required; set DAY_PLANNER_DB_HOST, DAY_PLANNER_DB_USER, DAY_PLANNER_DB_PASSWORD, and DAY_PLANNER_DB_NAME")

func Open(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		var err error
		databaseURL, err = URLFromEnv()
		if err != nil {
			return nil, err
		}
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

func URLFromEnv() (string, error) {
	host := os.Getenv("DAY_PLANNER_DB_HOST")
	user := os.Getenv("DAY_PLANNER_DB_USER")
	password := os.Getenv("DAY_PLANNER_DB_PASSWORD")
	name := os.Getenv("DAY_PLANNER_DB_NAME")
	if host == "" || user == "" || password == "" || name == "" {
		return "", ErrDatabaseConfigRequired
	}

	port := os.Getenv("DAY_PLANNER_DB_PORT")
	if port == "" {
		port = "5432"
	}
	sslMode := os.Getenv("DAY_PLANNER_DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   name,
	}
	query := u.Query()
	query.Set("sslmode", sslMode)
	u.RawQuery = query.Encode()
	return u.String(), nil
}
