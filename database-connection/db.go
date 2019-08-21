package main

// Dependencies
// MySQL Driver: go get -u github.com/go-sql-driver/mysql
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	dataSourceName = "root:root@tcp(localhost:3012)/golang_playground"
)

type User struct {
	Id int
	FirstName string
	LastName string
	EmailAddress string
	NullValue sql.NullInt64
}

func getDb() *sql.DB {
	db, err := sql.Open("mysql", dataSourceName)
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

func getUsers(db *sql.DB, users chan User) {
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

func createUser(db *sql.DB) {
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