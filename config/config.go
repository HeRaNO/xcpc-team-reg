package config

import (
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

var conf *Configure

type Configure struct {
	RDB     *RDBConfig     `yaml:"RDB"`
	Redis   *RedisConfig   `yaml:"Redis"`
	Srv     *SrvConfig     `yaml:"Server"`
	Contest *ContestConfig `yaml:"Contest"`
	Const   *ConstConfig   `yaml:"Const"`
}

type RDBConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
	TimeZone string `yaml:"TimeZone"`
}

type RedisConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
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

type ContestConfig struct {
	Name      string `yaml:"Name"`
	StartTime string `yaml:"StartTime"`
	EndTime   string `yaml:"EndTime"`
	Note      string `yaml:"Note"`
}

type ConstConfig struct {
	MaxTeamMember     int32    `yaml:"MaxTeamMember"`
	MaxTeamNameLength int      `yaml:"MaxTeamNameLength"`
	UserTokenLength   int      `yaml:"UserTokenLength"`
	JWTSecret         string   `yaml:"JWTSecret"`
	ValidStuIDLength  []int    `yaml:"ValidStuIDLength"`
	SchoolName        []string `yaml:"SchoolName"`
}

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

func InitConfig(file *string) {
	initConfigFile(file)

	var wg sync.WaitGroup
	wg.Add(1)
	go initDb(&wg)
	wg.Add(1)
	go initConst(&wg)
	wg.Add(1)
	go initServer(&wg)
	wg.Add(1)
	go initContest(&wg)
	initRedis()

	wg.Wait()
	log.Println("[INFO] init finished successfully")
}
