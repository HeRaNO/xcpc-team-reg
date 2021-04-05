package config

import (
	"log"
	"sync"
	"time"
)

const (
	FingerPrintTokenLength = 20
	JWTIDName              = "id"
	JWTAdminName           = "ro"
	JWTFingerPrintName     = "_f"
)

type CtxUserInfo struct {
	ID    int64
	Admin bool
}

type CtxName string

const (
	CtxUserInfoName CtxName = "user"
	MaxUploadSize   int64   = 2 << 20
)

var SchoolMap map[int]string
var StuIDMap map[int]bool
var MaxTeamNameLength, UserTokenLength int
var MaxTeamMember int32
var JWTSecret []byte

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
	MaxTeamMember = config.MaxTeamMember
	MaxTeamNameLength, UserTokenLength = config.MaxTeamNameLength, config.UserTokenLength
	JWTSecret = []byte(config.JWTSecret)
	log.Println("[INFO] init const finished successfully")
}
