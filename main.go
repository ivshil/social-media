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
	ID        int       `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	BirthDate time.Time `db:"birth_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Friend struct {
	ID            int
	InitiatorUser User
	SecondUser    User
	Status        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
type FriendDTO struct {
	ID              int       `db:"id"`
	InitiatorUserID int       `db:"initiator_user_id"`
	SecondUserID    int       `db:"second_user_id"`
	Status          int       `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
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
	// db.MustExec(schema)

	// userOne := User{
	// 	FirstName: "Ivan",
	// 	LastName:  "Shishman",
	// 	Email:     "ivsh@sh.c",
	// 	BirthDate: time.Now(),
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }

	// userTwo := User{
	// 	FirstName: "Grozdan",
	// 	LastName:  "Cvetkov",
	// 	Email:     "g.cvetkov@pete.bg",
	// 	BirthDate: time.Now(),
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }
	// query := `INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at)
	// VALUES (:first_name, :last_name, :email, :birth_date, :created_at, :updated_at)`

	// tx := db.MustBegin()
	// tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "Jason", "Moiron", "jmoiron@jmoiron.net", time.Now(), time.Now(), time.Now())
	// tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "John", "Doe", "johndoeDNE@gmail.net", time.Now(), time.Now(), time.Now())
	// tx.NamedExec(query, userOne)
	// tx.NamedExec(query, userTwo)
	// tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "Zuck", "ZeCuck", "zuki@fb.net", time.Now(), time.Now(), time.Now())
	// tx.MustExec("INSERT INTO users (first_name, last_name, email, birth_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", "Richard", "Brandson", "fendde@xyz.xyz", time.Now(), time.Now(), time.Now())
	// tx.MustExec("INSERT INTO friends (initiator_user_id, second_user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
	// 	1, 4, 3, time.Now(), time.Now())
	// tx.MustExec("INSERT INTO friends (initiator_user_id, second_user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
	// 	1, 2, 3, time.Now(), time.Now())
	// tx.MustExec("INSERT INTO friends (initiator_user_id, second_user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
	// 	1, 3, 1, time.Now(), time.Now())
	// tx.MustExec("INSERT INTO friends (initiator_user_id, second_user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
	// 	2, 4, 2, time.Now(), time.Now())
	// tx.MustExec("INSERT INTO friends (initiator_user_id, second_user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
	// 	3, 4, 3, time.Now(), time.Now())

	// tx.Commit()

	allUsers, err := GetAllUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range allUsers {
		fmt.Printf("User ID: %d, Name: %s %s, Email: %s\n",
			user.ID, user.FirstName, user.LastName, user.Email)
	}

	friendsDTO, err := GetFriendsForUser(db, 1) // Change the userID to the desired user's ID
	if err != nil {
		log.Fatal(err)
	}

	// Convert FriendDTO list to Friend list
	friends, err := ConvertFriendDTOToFriend(db, friendsDTO)
	if err != nil {
		log.Fatal(err)
	}

	// Print the converted Friend list
	for _, friend := range friends {
		status := ""
		if friend.Status == 1 {
			status = "pending"
		} else if friend.Status == 2 {
			status = "rejected"
		} else {
			status = "friends"
		}
		fmt.Printf("Initiator: %s %s, Second User: %s %s, Friend Status: %s\n",
			friend.InitiatorUser.FirstName, friend.InitiatorUser.LastName,
			friend.SecondUser.FirstName, friend.SecondUser.LastName, status)
	}
}

func GetFriendsForUser(db *sqlx.DB, userID int) ([]FriendDTO, error) {
	var friends []FriendDTO
	query := `
		SELECT f.id, f.initiator_user_id, f.second_user_id, f.status, f.created_at, f.updated_at
    	FROM friends f
    	WHERE f.initiator_user_id = $1 OR f.second_user_id = $1
    `
	err := db.Select(&friends, query, userID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func GetAllFriends(db *sqlx.DB, userID int) ([]Friend, error) {
	var friends []Friend
	query := `
        SELECT f.*, u1.*, u2.*
        FROM friends f
        JOIN users u1 ON f.initiator_user_id = u1.id
        JOIN users u2 ON f.second_user_id = u2.id
        WHERE f.initiator_user_id = $1 OR f.second_user_id = $1
    `
	err := db.Select(&friends, query, userID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	query := "SELECT * FROM users"
	err := db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(db *sqlx.DB, userID int) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE id = $1"
	err := db.Get(&user, query, userID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func ConvertFriendDTOToFriend(db *sqlx.DB, friendsDTO []FriendDTO) ([]Friend, error) {
	var friends []Friend

	for _, friendDTO := range friendsDTO {
		initiatorUser, err := GetUserByID(db, friendDTO.InitiatorUserID)
		if err != nil {
			return nil, err
		}

		secondUser, err := GetUserByID(db, friendDTO.SecondUserID)
		if err != nil {
			return nil, err
		}

		friend := Friend{
			ID:            friendDTO.ID,
			InitiatorUser: initiatorUser,
			SecondUser:    secondUser,
			Status:        friendDTO.Status,
			CreatedAt:     friendDTO.CreatedAt,
			UpdatedAt:     friendDTO.UpdatedAt,
		}

		friends = append(friends, friend)
	}

	return friends, nil
}
