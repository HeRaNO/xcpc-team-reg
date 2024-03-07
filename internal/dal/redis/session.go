package redis

import (
	"context"
	"errors"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

func GetSession(ctx context.Context, sid *string) (int64, error) {
	key := makeSessionKey(sid)
	ret, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		hlog.Infof("GetSession(): key is nil, sid: %d", sid)
		return 0, nil
	} else if err != nil {
		hlog.Errorf("GetSession(): redis query error, err: %+v", err)
		return 0, err
	}
	uid, err := strconv.ParseInt(ret, 10, 64)
	if err != nil {
		hlog.Errorf("GetSession(): parse user id error, err: %+v", err)
		return 0, err
	}
	return uid, nil
}

func SetSession(ctx context.Context, sid *string, uid int64) error {
	key := makeSessionKey(sid)
	_, err := redisClient.Set(ctx, key, uid, 0).Result()
	if err != nil {
		hlog.Errorf("SetSession(): redis set error, err: %+v", err)
		return err
	}
	return nil
}

func DelSession(ctx context.Context, sid *string) error {
	key := makeSessionKey(sid)
	n, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		hlog.Errorf("DelSession(): redis del error, err: %+v", err)
		return err
	}
	if n == 0 {
		err := errors.New("no sid found")
		hlog.Errorf("DelSession(): redis del error, err: %+v", err)
		return err
	}
	return nil
}
