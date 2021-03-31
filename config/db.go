package config

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var RDB *gorm.DB

func initDb(wg *sync.WaitGroup) {
	defer wg.Done()

	var err error
	config := conf.RDB
	if config == nil {
		panic("[FAILED] config file failed - RDB")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		config.Host, config.Username, config.Password, config.DBName, config.Port, config.TimeZone)
	RDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("[FAILED] init RDB failed")
		panic(err)
	}
	log.Println("[INFO] init database finished successfully")
}
