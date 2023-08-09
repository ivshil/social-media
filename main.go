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
	id serial PRIMARY KEY,
    first_name varchar(50) NOT NULL,
    last_name varchar(50),
    email varchar(50) NOT NULL,
	birth_date timestamp,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);

CREATE TABLE friends (
    id serial PRIMARY KEY,
    initiator_user_id int REFERENCES users(id) NOT NULL,
    second_user_id int REFERENCES users(id) NOT NULL,
    status integer NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE chats (
    id serial PRIMARY KEY,
    user_owner_id int REFERENCES users(id) NOT NULL,
    status boolean NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE chat_members (
    id serial PRIMARY KEY,
    chat_id int REFERENCES chats(id) NOT NULL,
    user_id int REFERENCES users(id) NOT NULL,
    join_date date NOT NULL,
    status varchar(10) NOT NULL,
    UNIQUE (chat_id, user_id)
);

CREATE TABLE messages (
    id serial PRIMARY KEY,
    created_at timestamp NOT NULL,
    chat_id int REFERENCES chats(id) NOT NULL,
    user_id int REFERENCES users(id) NOT NULL,
    message_content_link varchar,
    preview varchar
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

	userOne := User{
		FirstName: "Ivan",
		LastName:  "Shishman",
		Email:     "ivsh@sh.c",
		BirthDate: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userTwo := User{
		FirstName: "Grozdan",
		LastName:  "Cvetkov",
		Email:     "g.cvetkov@pete.bg",
		BirthDate: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	query := `INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at)
	VALUES (:first_name, :last_name, :email, :birth_date, :created_at, :updated_at)`

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "Jason", "Moiron", "jmoiron@jmoiron.net", time.Now(), time.Now(), time.Now())
	tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "John", "Doe", "johndoeDNE@gmail.net", time.Now(), time.Now(), time.Now())
	tx.NamedExec(query, userOne)
	tx.NamedExec(query, userTwo)
	tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "Zuck", "ZeCuck", "zuki@fb.net", time.Now(), time.Now(), time.Now())
	tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "Richard", "Brandson", "fendde@xyz.xyz", time.Now(), time.Now(), time.Now())
	tx.Commit()

	/

	users := []User{}
	db.Select(&users, "SELECT * FROM person ORDER BY first_name ASC")
	first, second := users[0], users[1]

	fmt.Printf("\n%s, %s\n %s, %s, \n", first.FirstName, first.CreatedAt, second.FirstName, second.CreatedAt)
}
