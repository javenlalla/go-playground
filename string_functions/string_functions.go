package main

import (
	"log"
	"strconv"
	"strings"
)

func main() {
	log.Println(ConvertStringToFloat("129.0000"))
	log.Println(ConvertStringToInt("365"))
}

// ConvertStringToFloat converts a number stored as a string, such as "69.0000", and
// converts it to a float and returns the updated type value.
func ConvertStringToFloat(p string) float64 {
	if p == "" {
		return 0
	}

	n, err := strconv.ParseFloat(p, 64)
	if err != nil {
		log.Fatalf("Unable to convert numbered stored as string `%s` to float: %s", p, err)
	}

	return n
}

// ConvertStringToInt converts a number stored as a string, such as "365", and
// converts it to a int and returns the updated type value.
func ConvertStringToInt(p string) int64 {
	if p == "" {
		return 0
	}

	n, err := strconv.ParseInt(p, 10,64)
	if err != nil {
		log.Fatalf("Unable to convert numbered stored as string `%s` to int: %s", p, err)
	}

	return n
}

// ConvertStringBoolToIntFlag converts a boolean value stored as a string, such as "True" or "False", and
// converts it to an int denoting `1` for `true` and `0` for `false` and returns the number value.
func ConvertStringBoolToIntFlag(p string) int64 {
	p = strings.ToLower(p)

	if p == "true" {
		return 1
	}

	return 0
}
