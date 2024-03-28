package rdb

import (
	"fmt"

	"github.com/HeRaNO/xcpc-team-reg/internal/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init(conf *config.RDBConfig) {
	if conf == nil {
		hlog.Fatal("RDB config failed: conf is nil")
		panic("make static check happy")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		conf.Host, conf.Username, conf.Password, "xcpc_team_reg", conf.Port)
	RDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		hlog.Fatalf("init RDB failed, err: %+v", err)
	}
	db = RDB
	hlog.Info("init database finished successfully")
}
