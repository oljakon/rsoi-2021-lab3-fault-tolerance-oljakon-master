package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

//const dsn = "postgresql://program:test@localhost:5432/services?sslmode=disable"
const dsn = "postgres://scbufadmtlaini:d9905cfb6de475729ec7753b5c7d5ac3c31d9f93cc0a4163e0da19af9b132ace@ec2-174-129-16-183.compute-1.amazonaws.com:5432/d4ah6bmp8va57l"

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
