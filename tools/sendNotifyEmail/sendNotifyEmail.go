package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"text/template"

	"github.com/jordan-wright/email"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Configure struct {
	RDB *RDBConfig `yaml:"RDB"`
	Srv *SrvConfig `yaml:"Server"`
}

type RDBConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
	TimeZone string `yaml:"TimeZone"`
}

type SrvConfig struct {
	Domain        string            `yaml:"Domain"`
	Port          int               `yaml:"Port"`
	SMTPAddr      string            `yaml:"SMTPAddr"`
	SMTPPort      int               `yaml:"SMTPPort"`
	EmailAddr     string            `yaml:"EmailAddr"`
	EmailPassword string            `yaml:"EmailPassword"`
	EmailServer   string            `yaml:"EmailServer"`
	EmailSign     string            `yaml:"EmailSign"`
	EmailAlias    string            `yaml:"EmailAlias"`
	EmailAction   map[string]string `yaml:"EmailAction"`
	EmailSubject  map[string]string `yaml:"EmailSubject"`
}

var conf Configure

func initConfigFile(filePath *string) {
	fileBytes, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Println("[FAILED] read config file failed")
		panic(err)
	}
	if err = yaml.Unmarshal(fileBytes, &conf); err != nil {
		log.Println("[FAILED] unmarshal yaml file failed")
		panic(err)
	}
}

var RDB *gorm.DB

func initDb() {
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

var Domain, ListenPort string
var SMTPAddr, SMTPHost string
var SMTPPort int
var EmailSign, EmailAddr, EmailPassword, EmailServer, EmailAlias, EmailFrom string
var EmailActionMap map[string]string
var EmailSubjectMap map[string]string

func initEmail() {
	if conf.Srv == nil {
		panic("[FAILED] config file error - Email")
	}
	config := conf.Srv
	Domain = config.Domain
	ListenPort = fmt.Sprintf(":%d", config.Port)
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
	log.Println("[INFO] init email finished successfully")
}

func initConfig(filePath *string) {
	initConfigFile(filePath)
	initDb()
	initEmail()
}

const (
	TableUserInfo = "t_user"
)

type User struct {
	UserID     int64  `gorm:"column:user_id;primaryKey" json:"userid"`
	Name       string `gorm:"column:user_name" json:"name"`
	Email      string `gorm:"column:email" json:"email"`
	School     int    `gorm:"column:school" json:"school"`
	StuID      string `gorm:"column:stu_id" json:"stuid"`
	BelongTeam int64  `gorm:"column:belong_team" json:"teamid"`
	IsAdmin    int    `gorm:"column:is_admin" json:"is_admin"`
}

func sendEmail(name *string, emailRecv *string) error {
	tmpl, err := template.ParseFiles("./email_template.tmpl")

	if err != nil {
		return err
	}

	content := new(bytes.Buffer)
	tmpl.Execute(content, struct {
		Name        string
		ContestName string
		EndTime     string
		Sign        string
	}{
		Name:        *name,
		ContestName: "电子科技大学第十九届程序设计竞赛（补报）",
		EndTime:     "2021-05-11 11:59:59",
		Sign:        EmailSign,
	})

	e := &email.Email{
		To:      []string{*emailRecv},
		From:    EmailFrom,
		Subject: "校赛报名通知",
		HTML:    content.Bytes(),
	}

	auth := smtp.PlainAuth("", EmailAddr, EmailPassword, EmailServer)
	err = e.SendWithTLS(SMTPHost, auth, &tls.Config{ServerName: EmailServer})

	return err
}

func main() {
	configFilePath := flag.String("c", "./conf/config.yaml", "the path of configure file")
	flag.Parse()

	initConfig(configFilePath)

	rec := make([]User, 0)
	result := RDB.Model(&User{}).Table(TableUserInfo).Where("belong_team = ?", 0).Find(&rec)

	if result.Error != nil {
		panic(result.Error)
	}

	failed := make([]int64, 0)

	for _, usr := range rec {
		err := sendEmail(&usr.Name, &usr.Email)
		if err != nil {
			log.Printf("[ERROR] send email error, usr: %+v, err: %s", usr, err.Error())
			failed = append(failed, usr.UserID)
		} else {
			log.Printf("[ERROR] send email ok, user_id: %d", usr.UserID)
		}
	}

	log.Printf("[INFO] send email finished. failed id: %+v", failed)
}
