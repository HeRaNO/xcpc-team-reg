package config

import (
	"log"
	"sync"
)

var SchoolMap map[int]string
var StuIDMap map[int]bool

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
	log.Println("[INFO] init const finished successfully")
}
