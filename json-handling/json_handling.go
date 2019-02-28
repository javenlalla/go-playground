package main

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
	var jsonBlob = []byte(`[
		{"name": "Test Name", "email": "test1@test.com"},
		{"name": "Testing Name 2",    "email": "test2@test.com"}
	]`)

	var users []User
	err := json.Unmarshal(jsonBlob, &users)
	if err != nil {
		log.Fatalf("Unable to parse JSON: %s", err)
	}

	log.Println(users)
}
