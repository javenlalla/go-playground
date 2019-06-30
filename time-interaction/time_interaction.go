package main

import (
	"log"
	"time"
)

func main() {
	mysqlTimestamp := getCurrentTimeAsMysql()
	log.Println(mysqlTimestamp)

	date := "2018-07-05"
	dateObj := convertDateTimeStringToTime(date)

	log.Println(dateObj)

	timestamp := "2018-07-05 04:23:59"
	timeObj := convertDateTimeStringToTime(timestamp)
	location, _ := time.LoadLocation("Europe/Berlin")

	log.Println(timeObj, timeObj.In(location))
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

func convertDateTimeStringToTime(s string) time.Time {
	location, _ := time.LoadLocation("America/New_York")
	timeObj, err := time.ParseInLocation("2006-01-02 15:04:05", s, location)
	if err != nil {
		log.Fatalf("Unable to parse datetime string `%s`: %s", s, err)
	}

	return timeObj
}

func getCurrentTimeAsMysql() string {
	now := time.Now()

	return now.Format("2006-01-02 15:04:05")
}
