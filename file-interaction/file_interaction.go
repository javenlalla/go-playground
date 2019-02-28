package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	readingFiles()
	writingFiles()
}

func writingFiles() {
	filePath := "users_test.csv"
	file := openFileWithWrite(filePath)

	_, err := file.WriteString("test first name,another test value, thirdTestValue\n")
	if err != nil {
		log.Fatalf("Unable to write to file `%s`: %s", filePath, err)
	}
}

func openFileWithWrite(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Unable to open file `%s`: %s", filePath, err)
	}

	return file
}

func readingFiles() {
	filePath := "users.csv"
	file := openFile(filePath)

	userRows := make(chan string)
	go getUsersFromFile(file, userRows)

	for userRow := range userRows {
		log.Println(userRow)
	}
}

func openFile(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("Unable to open file `%s`: %s", filePath, err)
	}

	return file
}

func getUsersFromFile(file *os.File, userRows chan string) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		userRows <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file: %s", err)
	}

	close(userRows)
}
