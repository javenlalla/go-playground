package main

// Dependencies
//
// MySQL Driver:
// go get -u github.com/go-sql-driver/mysql
//
// sqlx: https://github.com/jmoiron/sqlx
// go get github.com/jmoiron/sqlx

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	dataSourceName = "root:root@tcp(localhost:3012)/golang_playground"
)

type User struct {
	Id int `db:"id"`
	FirstName string `db:"first_name"`
	LastName string `db:"last_name"`
	EmailAddress string `db:"email"`
	NullValue sql.NullInt64 `db:"null_column"`
}

func getDb() *sqlx.DB {
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	return db
}

func main() {
	db := getDb()

	log.Println("Fetching users via row marshalling.")
	log.Println(getUsersViaRowMarshal(db))

	log.Println("Create a User via row marshalling.")
	createUserViaRowMarshal(db, User{
		FirstName: "row",
		LastName: "marshal",
		EmailAddress: "row@marshal.com",
		NullValue: sql.NullInt64{},
	})

	log.Println("Fetching non-existent User via marshalling.")
	log.Println(getNonExistentUserViaRowMarshal(db))

	log.Println("Executing database operations using native sql API.")
	users := make(chan User)
	go getUsers(db, users)

	for user := range users {
		log.Println(user)
	}

	createUser(db)

	err := db.Close()
	if err != nil {
		log.Fatalf("Unable to close database connection: %s", err)
	}
}

func getUsersViaRowMarshal(db *sqlx.DB) []User {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
		    null_column
		FROM
			golang_playground.users
		WHERE
			id >= ?
	`

	var users []User
	err := db.Select(&users, query, 35)
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
			id,
			first_name,
			last_name,
			email,
		    null_column
		FROM
			golang_playground.users
		WHERE
			id = ?
	`

	var user User
	err := db.Get(&user, query, 100)
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
			id,
			first_name,
			last_name,
			email,
		    null_column
		FROM
			golang_playground.users
		WHERE
			id >= ?
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Unable to prepare query: %s. Error: %s", query, err)
	}

	rows, err := stmt.Query(35)
	if err != nil {
		log.Fatalf("Error executing query: %s. Error: %s", query, err)
	}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.EmailAddress, &user.NullValue); err != nil {
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
			golang_playground.users (first_name, last_name, email, null_column)
		VALUES (:first_name, :last_name, :email, :null_column)
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
			golang_playground.users (first_name, last_name, email, null_column)
		VALUES (?, ?, ?, ?)
	`

	stmt, err := db.Prepare(insert)
	if err != nil {
		log.Fatalf("Unable to prepare query: %s. Error: %s", insert, err)
	}

	res, err := stmt.Exec("test insert", "test insert last name", "test@test.com", sql.NullInt64{})
	if err != nil {
		log.Fatalf("Unable to execute query: %s. Error: %s", insert, err)
	}

	log.Println(res.LastInsertId())
	log.Println(res.RowsAffected())
}