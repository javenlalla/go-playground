package main

// Dependencies:
// BBolt DB: go get go.etcd.io/bbolt/...
import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"log"
	"strconv"
)

type User struct {
	Id int
	Name string
	EmailAddress string
}

func main() {
	dbPath := "mydb"
	db := getBoltDb(dbPath)
	defer db.Close()


	insertData(db)

	users := make(chan User)
	go readData(db, users)

	for user := range users {
		log.Println(user)
	}

	if _, e := getUserRecord(db, 123452); e == false {
		log.Println("Record not found")
	}

	u, e := getUserRecord(db, 1234)
	if e == false {
		log.Println("Record not found")
	} else {
		log.Println(u)
	}
}

// getUserRecord fetches a User record from the Users Bucket by the provided User ID.
func getUserRecord(db *bbolt.DB, id int) (User, bool) {
	var u User
	exists := false

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		v := b.Get(itob(id))
		if len(v) > 0 {
			err := json.Unmarshal(v, &u)
			if err != nil {
				log.Fatalf("Error parsing record: %s", err)
			}
			exists = true
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error retrieving record by key: %s", err)
	}

	return u, exists
}

// readData retrieves records from the Users Bucket.
func readData(db *bbolt.DB, users chan User) {
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u User
			err := json.Unmarshal(v, &u)
			if err != nil {
				log.Fatalf("Error parsing User record %s: %s", string(v), err)
			}

			users <- u
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Unable to retrieve users from Bucket: %s", err)
	}

	close(users)
}

// insertData inserts a set of users into a BBolt Bucket.
func insertData(db *bbolt.DB) {
	users := []User{
		{Id: 1234, Name: "Roger", EmailAddress: "rtyres2@imageshack.us"},
		{Id: 3498, Name: "Hirsch", EmailAddress: "hcatmulld@ox.ac.uk"},
	}

	err := db.Update(func(tx *bbolt.Tx) error {
		//b := getBucketForWriting(tx, "users")
		b := tx.Bucket([]byte("users"))

		for _, user := range users {
			userRecord, err := json.Marshal(user)
			if err != nil {
				return err
			}

			err = b.Put(itob(user.Id), userRecord)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Unable to create user: %s", err)
	}
}

// itob converts the provided int value to a slice of bytes.
func itob(i int) []byte {
	return []byte(strconv.Itoa(i))
}

// getBoltDb gets an BBolt DB instance at the provided path.
func getBoltDb(dbPath string) *bbolt.DB {
	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		log.Fatalf("unable to open BBolt DB file: %s", err)
	}

	setupDb(db)

	return db
}

// setupDb initializes all Buckets in the provided BBolt DB instance.
func setupDb(db *bbolt.DB) {
	buckets := []string{
		"users",
	}

	for _, b := range buckets {
		err := db.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				log.Fatalf("unable to open `%s` Bucket: %s", b, err)
			}

			return nil
		})
		if err != nil {
			log.Fatalf("Unable to create `%s` Bucket: %s", b, err)
		}
	}
}
