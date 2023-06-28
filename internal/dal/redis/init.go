package redis

import (
	"fmt"

	"github.com/HeRaNO/xcpc-team-reg/internal/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func Init(conf *config.RedisConfig) {
	if conf == nil {
		hlog.Fatal("Redis config failed: conf is nil")
		panic("make static check happy")
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password:     conf.Password,
		MaxIdleConns: maxIdle,
	})
	if redisClient == nil {
		hlog.Fatal("init Redis failed: redisClient == nil")
	}
	hlog.Info("init Redis finished successfully")
}
