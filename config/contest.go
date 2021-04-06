package config

import (
	"log"
	"sync"
	"time"
)

var ContestName, ContestNote string
var ContestStartTime, ContestEndTime time.Time

func initContest(wg *sync.WaitGroup) {
	defer wg.Done()

	var err error
	config := conf.Contest
	if config == nil {
		panic("[FAILED] config file failed - RDB")
	}
	ContestName, ContestNote = config.Name, config.Note
	startTime := config.StartTime
	loc, _ := time.LoadLocation("Local")
	ContestStartTime, err = time.ParseInLocation("2006-01-02 15:04:05", startTime, loc)
	if err != nil {
		log.Println("[FAILED] parse contest start time failed")
		panic(err)
	}
	endTime := config.EndTime
	ContestEndTime, err = time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
	if err != nil {
		log.Println("[FAILED] parse contest end time failed")
		panic(err)
	}
	if !ContestStartTime.Before(ContestEndTime) {
		log.Panicln("[FAILED] contest start time not before contest end time")
	}

	log.Println("[INFO] init contest finished successfully")
}
