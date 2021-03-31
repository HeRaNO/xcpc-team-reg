package config

import (
	"fmt"
	"log"
	"sync"
)

var ListenPort string
var EmailTokenLength int
var SMTPAddr, SMTPHost string
var SMTPPort int
var EmailSign, EmailAddr, EmailPassword, EmailServer, EmailAlias, EmailFrom string
var EmailActionMap map[string]string
var EmailSubjectMap map[string]string

func initServer(wg *sync.WaitGroup) {
	defer wg.Done()

	if conf.Srv == nil {
		panic("[FAILED] config file error - Server")
	}
	config := conf.Srv
	ListenPort = fmt.Sprintf(":%d", config.Port)
	EmailTokenLength = config.EmailTokenLength
	EmailSign, EmailAddr, EmailPassword = config.EmailSign, config.EmailAddr, config.EmailPassword
	EmailServer, EmailAlias = config.EmailServer, config.EmailAlias
	EmailFrom = fmt.Sprintf("%s <%s>", EmailAlias, EmailAddr)
	SMTPAddr = config.SMTPAddr
	SMTPPort = config.SMTPPort
	SMTPHost = fmt.Sprintf("%s:%d", SMTPAddr, SMTPPort)
	EmailActionMap = make(map[string]string)
	for k, v := range config.EmailAction {
		EmailActionMap[k] = v
	}
	EmailSubjectMap = make(map[string]string)
	for k, v := range config.EmailSubject {
		EmailSubjectMap[k] = v
	}
	log.Println("[INFO] init server finished successfully")
}
