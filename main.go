package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var schema = `CREATE TABLE users (
    first_name varchar(50),
    last_name varchar(50),
    email varchar(50),
	birth_date timestamp,
	created_at timestamp,
	updated_at timestamp
)`

var (
	dbHost = ""
	dbPort = ""
	dbUser = ""
	dbPass = ""
	dbName = ""
)

type User struct {
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	BirthDate time.Time `db:"birth_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbHost = os.Getenv("PGDB_HOST")
	dbPort = os.Getenv("PGDB_PORT")
	dbUser = os.Getenv("PGDB_USER")
	dbPass = os.Getenv("PGDB_PASS")
	dbName = os.Getenv("PGDB_NAME")

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPass, "127.0.0.1", dbPort, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, crreated_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "Jason", "Moiron", "jmoiron@jmoiron.net", time.Now(), time.Now(), time.Now())
	tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, crreated_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "John", "Doe", "johndoeDNE@gmail.net", time.Now(), time.Now(), time.Now())
	tx.Commit()

	users := []User{}
	db.Select(&users, "SELECT * FROM person ORDER BY first_name ASC")
	first, second := users[0], users[1]

	fmt.Printf("\n%s, %s\n %s, %s, \n", first.FirstName, first.CreatedAt, second.FirstName, second.CreatedAt)
}
