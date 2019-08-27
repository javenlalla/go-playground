package main

import (
	"log"
	"time"
)

func main() {
	mysqlTimestamp := getCurrentTimeAsMysql()
	log.Println(mysqlTimestamp)

	date := "2018-07-05"
	dateObj := convertDateStringToTime(date)

	log.Println(dateObj)

	timestamp := "2018-07-05 04:23:59"
	timeObj := convertDateTimeStringToTime(timestamp)
	location, _ := time.LoadLocation("Europe/Berlin")

	log.Println(timeObj, timeObj.In(location))

	date = "2018-05-15"
	dateObj = addDaysToDate(date, 5)
	log.Println(dateObj)
	dateObj = subtractDaysFromDate(date, 10)
	log.Println(dateObj)
}

// convertDateStringToTime converts a string formatted as YYYY-MM-DD to a time.Time instance.
// Note: because the a time string is not included, it is not necessary to include location or parsing in location.
func convertDateStringToTime(s string) time.Time {
	timeObj, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalf("Unable to parse date string `%s`: %s", s, err)
	}

	return timeObj
}

// convertDateTimeStringToTime converts a string formatted as YYYY-MM-DD hh:mm:ss to a time.Time instance.
// Note: the location is used under the assumption that it is known the provided time string is UTC.
func convertDateTimeStringToTime(s string) time.Time {
	location, _ := time.LoadLocation("America/New_York")
	timeObj, err := time.ParseInLocation("2006-01-02 15:04:05", s, location)
	if err != nil {
		log.Fatalf("Unable to parse datetime string `%s`: %s", s, err)
	}

	return timeObj
}

// getCurrentTimeAsMysql gets the current time as a time.Time instance returns as a
// MySQL formatted DATETIME.
func getCurrentTimeAsMysql() string {
	now := time.Now()

	return now.Format("2006-01-02 15:04:05")
}

// subtractDaysFromDate subtracts the provided `n` number of days from the provided date string
//// and returns the updated time.Time instance.
func subtractDaysFromDate(s string, n int) time.Time {
	timeObj, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalf("Unable to parse date string `%s`: %s", s, err)
	}

	updatedTimeObj := timeObj.Add(time.Hour * 24 * time.Duration(- n))

	return updatedTimeObj
}

// addDaysToDate adds the provided `n` number of days to the provided date string
// and returns the updated time.Time instance.
func addDaysToDate(s string, n int64) time.Time {
	timeObj, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalf("Unable to parse date string `%s`: %s", s, err)
	}

	updatedTimeObj := timeObj.Add(time.Hour * 24 * time.Duration(n))

	return updatedTimeObj
}
