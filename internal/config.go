package internal

import (
	"fmt"
	"os"

	"github.com/HeRaNO/xcpc-team-reg/internal/config"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/internal/email"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/goccy/go-yaml"
)

var Domain string

func initConfigFile(filePath *string) *config.Configure {
	fileBytes, err := os.ReadFile(*filePath)
	if err != nil {
		hlog.Fatalf("read config file failed, err: %+v", err)
	}
	conf := config.Configure{}
	if err = yaml.Unmarshal(fileBytes, &conf); err != nil {
		hlog.Fatalf("unmarshal yaml file failed, err: %+v", err)
	}
	return &conf
}

func initSrv(conf *config.SrvConfig) string {
	if conf == nil {
		hlog.Fatal("Srv config failed: conf is nil")
		panic("make static check happy")
	}
	Domain = conf.Domain
	return fmt.Sprintf(":%d", conf.Port)
}

func InitConfig(filePath *string) string {
	conf := initConfigFile(filePath)
	rdb.Init(conf.RDB)
	redis.Init(conf.Redis)
	email.Init(conf.Email)
	contest.Init(conf.Contest)

	return initSrv(conf.Srv)
}
