package config

import (
	"log"
	"sync"
	"time"
)

var SchoolMap map[int]string
var StuIDMap map[int]bool
var MaxTeamNameLength, UserTokenLength int

const (
	LOGIN_EXPIRETIME      = 24 * time.Hour
	EMAILTOKEN_EXPIRETIME = 10 * time.Minute
	EMAILSEND_GAPTIME     = 2 * time.Minute
)

func initConst(wg *sync.WaitGroup) {
	defer wg.Done()

	config := conf.Const
	if config == nil {
		panic("[FAILED] config file failed - RDB")
	}
	SchoolMap = make(map[int]string)
	StuIDMap = make(map[int]bool)
	for i, schoolName := range config.SchoolName {
		SchoolMap[i+1] = schoolName
	}
	for _, stuIDLength := range config.ValidStuIDLength {
		StuIDMap[stuIDLength] = true
	}
	MaxTeamNameLength, UserTokenLength = config.MaxTeamNameLength, config.UserTokenLength
	log.Println("[INFO] init const finished successfully")
}
