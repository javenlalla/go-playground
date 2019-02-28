package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type User struct {
	ID     int  `json:"id"`
	Name   string `json:"name"`
	EmailAddress string `json:"email"`
}

type ResponseBody struct {
	ID     int  `json:"id"`
	Name   string `json:"name"`
	EmailAddress string `json:"email"`
}

func main() {
	hc := getHttpClient()

	getCall(hc)
	postCall(hc)
}

func getHttpClient() *http.Client {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 15 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 15 * time.Second,
	}

	var httpClient = &http.Client{
		Timeout: time.Second * 15,
		Transport: netTransport,
	}

	return httpClient
}

func postCall(hc *http.Client) {
	user := User{
		ID: 12345,
		Name: "Test Post Request",
		EmailAddress: "test@test.com",
	}

	var responseBody ResponseBody
	url := 	"https://jsonplaceholder.typicode.com/users"

	err := executePostCall(hc, url, user, &responseBody)
	if err != nil {
		log.Fatalf("Error executing POST call: %s", err)
	}

	log.Println(responseBody)
}

func getCall(hc *http.Client) {
	url := 	"https://jsonplaceholder.typicode.com/users"

	var users []User
	err := executeGetCall(hc, url, &users)
	if err != nil {
		log.Fatalf("Error executing GET call: %s", err)
	}

	log.Println(users)
}

// Update responseBody type
func executeGetCall(hc *http.Client, url string, responseBody *[]User) error {
	response, err := hc.Get(url)
	if err != nil {
		return fmt.Errorf("unable to execute http call against URL %s: %s", url, err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("API call did not return 200 for http call. Received: %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("unable to parse JSON Response for http call: %s", err)
	}

	return nil
}

// Update requestBody and responseBody types
func executePostCall(hc *http.Client, url string, requestBody User, responseBody *ResponseBody) error {
	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("unable to convert struct to JSON for http call: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return fmt.Errorf("unable to create request for http call: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		return fmt.Errorf("unable to execute http call: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("API call did not return 201 for http call. Received: %d", resp.StatusCode)
	}

	//getJsonRawOutput(resp)

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("unable to parse JSON Response for http call: %s", err)
	}

	return nil
}

// getJsonRawOutput prints the raw JSON string in the provided response
// This function is mainly used for debugging purposes and should NOT be used for parsing a JSON response.
// Note: when this function is called, the contents of the Body are no longer accessible for retrieval afterwards.
// In other words, this function can't be called and then another reader/parser function called on the same response
// immediately afterwards.
func getJsonRawOutput(response *http.Response) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Unable to read JSON response: %s", err)
	}

	log.Println(string(body))
}