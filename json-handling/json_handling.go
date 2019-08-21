package main

// General tips:
//
// 1. If the JSON data being unmarshaled is an array, create a Struct to represent each element and treat the target
// variable as a slice of said Struct. See `sampleJSONTypeOne`

import (
	"encoding/json"
	"log"
)

type User struct {
	ID     int  `json:"id"`
	Name   string `json:"name"`
	EmailAddress string `json:"email"`
	Pets []string `json:"pets"`
}

type UsersReport struct {
	Users []User `json:"UsersReport"`
}

const (
	sampleJSONTypeOne = `[
		{"name": "Test Name", "email": "test1@test.com"},
		{"name": "Testing Name 2",    "email": "test2@test.com"}
	]`

	sampleJSONTypeTwo = `{"UsersReport": [
		{"name": "Test Name", "email": "test1@test.com"},
		{"name": "Testing Name 2",    "email": "test2@test.com"}
	]}`
)

func main() {
	convertStructToJson()
	convertJsonToStruct()
}

func convertStructToJson() {
	user := User{
		ID:     1,
		Name:   "TestTest",
		EmailAddress: "test@test.com",
		Pets: []string{"dog", "cat", "bunny", "hamster"},
	}

	u, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Unable to convert struct to JSON: %s", err)
	}

	log.Println(string(u))
}

func convertJsonToStruct() {
	jsonBlob := []byte(sampleJSONTypeOne)

	var users []User
	err := json.Unmarshal(jsonBlob, &users)
	if err != nil {
		log.Fatalf("Unable to parse JSON: %s", err)
	}

	log.Println(users)

	jsonBlob = []byte(sampleJSONTypeTwo)

	var usersTwo UsersReport
	err = json.Unmarshal(jsonBlob, &usersTwo)
	if err != nil {
		log.Fatalf("Unable to parse JSON: %s", err)
	}

	log.Println(usersTwo)
}
