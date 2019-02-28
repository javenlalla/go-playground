package main

import (
	"log"
	"time"
)

func main() {
	mysqlTimestamp := getCurrentTimeAsMysql()
	log.Println(mysqlTimestamp)

	timestamp := "2018-07-05 04:23:59"
	timeObj := convertStringToTime(timestamp)
	location, _ := time.LoadLocation("Europe/Berlin")

	log.Println(timeObj, timeObj.In(location))
}

func convertStringToTime(s string) time.Time {
	location, _ := time.LoadLocation("America/New_York")
	timeObj, err := time.ParseInLocation("2006-02-01 15:04:05", s, location)
	if err != nil {
		log.Fatalf("Unable to parse time string: %s", err)
	}

	return timeObj
}

func getCurrentTimeAsMysql() string {
	now := time.Now()

	return now.Format("2006-02-01 15:04:05")
}
