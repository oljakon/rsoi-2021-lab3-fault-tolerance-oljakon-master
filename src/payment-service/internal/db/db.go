package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

//const dsn = "postgresql://program:test@localhost:5432/services?sslmode=disable"
const dsn = "postgres://xypzmkdvbismsy:999df7c4a9ad21ca8ec198940f8aac70b6e7547ed99f1bce5ade8a42812325e8@ec2-174-129-16-183.compute-1.amazonaws.com:5432/da63d9ck8m3u2q"

// CreateConnection to persons db
func CreateConnection() *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
