package main

// Dependencies
//
// Sqlite Driver (without CGO):
// go get modernc.org/sqlite
//
// sqlx: https://github.com/jmoiron/sqlx
// go get github.com/jmoiron/sqlx

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	_ "modernc.org/sqlite"
	"time"
)

const (
	dataSourceName = "root:root@tcp(localhost:3012)/golang_playground"
)

type User struct {
	RowID        int64         `db:"rowid"`
	FirstName    string        `db:"first_name"`
	LastName     string        `db:"last_name"`
	EmailAddress string        `db:"email"`
	NullValue    sql.NullInt64 `db:"null_column"`
	Created      sql.NullTime  `db:"created"`
	Modified     sql.NullTime  `db:"modified"`
}

func getDb() *sqlx.DB {
	db, err := sqlx.Open("sqlite", "users.db")
	if err != nil {
		log.Fatalf("unable to open database: %s", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}

	initializeTable(db)

	return db
}

func initializeTable(db *sqlx.DB) {
	q := `
		CREATE TABLE IF NOT EXISTS users (
		    first_name varchar(100) NOT NULL,
    		last_name varchar(100) NOT NULL,
    		email varchar(100) NOT NULL,
    		null_column varchar(50) DEFAULT NULL,
			created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    		modified datetime NOT NULL,
    		UNIQUE(email)
		);
	`

	_, err := db.Exec(q)
	if err != nil {
		log.Fatalf("error creating users table: %s", err)
	}
}

func main() {
	db := getDb()
	defer db.Close()

	log.Println("Create a User via row marshalling.")
	createUserViaRowMarshal(db, User{
		FirstName:    "first user",
		LastName:     "first last name",
		EmailAddress: "row@marshal.com",
		NullValue:    sql.NullInt64{},
		Modified:     sql.NullTime{Time: time.Now(), Valid: true},
		//Modified: "random string",
	})

	log.Println("Fetching users via row marshalling.")
	log.Println(getUsersViaRowMarshal(db))

	log.Println("Fetching non-existent User via marshalling.")
	log.Println(getNonExistentUserViaRowMarshal(db))

	log.Println("Executing database operations using native sql API.")
	users := make(chan User)
	go getUsers(db, users)

	for user := range users {
		log.Println(user)
	}

	createUser(db)
}

func getUsersViaRowMarshal(db *sqlx.DB) []User {
	query := `
		SELECT
			rowid,
			first_name,
			last_name,
			email,
		    null_column
		FROM
			users
		WHERE
			rowid = ?
	`

	var users []User
	err := db.Select(&users, query, 1)
	if err != nil {
		log.Fatalf("Unable to execute Select query: %s. Error: %s", query, err)
	}

	// Or if needed to load each User individually, use the snippet below instead
	//var users []User
	//rows, err := db.Queryx(query, 35)
	//if err != nil {
	//	log.Fatalf("Unable to execute query: %s. Error: %s", query, err)
	//}
	//for rows.Next() {
	//	user := User{}
	//	err := rows.StructScan(&user)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	users = append(users, user)
	//}

	return users
}

func getNonExistentUserViaRowMarshal(db *sqlx.DB) (User, bool) {
	query := `
		SELECT
			rowid,
			first_name,
			last_name,
			email,
		    null_column
		FROM
			users
		WHERE
			rowid = ?
	`

	var user User
	err := db.Get(&user, query, 10)
	if err == sql.ErrNoRows {
		return user, false
	} else if err != nil {
		log.Fatalf("Unable to execute Select query: %s. Error: %s", query, err)
	}

	return user, true
}

func getUsers(db *sqlx.DB, users chan User) {
	query := `
		SELECT
			rowid,
			first_name,
			last_name,
			email,
		    null_column,
		    created,
		    modified
		FROM
			users
		WHERE
			rowid >= ?
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Unable to prepare query: %s. Error: %s", query, err)
	}

	rows, err := stmt.Query(0)
	if err != nil {
		log.Fatalf("Error executing query: %s. Error: %s", query, err)
	}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.RowID, &user.FirstName, &user.LastName, &user.EmailAddress, &user.NullValue, &user.Created, &user.Modified); err != nil {
			log.Fatalf("Unable to scan row result: %s", err)
		}

		users <- user
	}

	close(users)

	err = rows.Close()
	if err != nil {
		log.Printf("Unable to close Rows after query result. Error: %s", err)
	}
}

func createUserViaRowMarshal(db *sqlx.DB, u User) {
	insert := `
		INSERT INTO
			users (first_name, last_name, email, null_column, modified)
		VALUES (:first_name, :last_name, :email, :null_column, :modified)
	`

	res, err := db.NamedExec(insert, u)
	if err != nil {
		log.Fatalf("Unable to execute Insert query: %s. Error: %s", insert, err)
	}

	log.Println(res.LastInsertId())
	log.Println(res.RowsAffected())
}

func createUser(db *sqlx.DB) {
	insert := `
		INSERT INTO
			users (first_name, last_name, email, null_column, modified)
		VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := db.Prepare(insert)
	if err != nil {
		log.Fatalf("Unable to prepare query: %s. Error: %s", insert, err)
	}

	res, err := stmt.Exec("test insert", "test insert last name", "test2@test.com", sql.NullInt64{}, sql.NullTime{time.Now(), true})
	if err != nil {
		log.Fatalf("Unable to execute query: %s. Error: %s", insert, err)
	}

	insertId, _ := res.LastInsertId()
	affectedRows, _ := res.RowsAffected()

	log.Println(fmt.Sprintf("Insert ID: %d", insertId))
	log.Println(fmt.Sprintf("Affected Rows: %d", affectedRows))
}
