package redis

import (
	"context"
	"strconv"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

func GetSession(ctx context.Context, sid *string) (int64, berrors.Berror) {
	key := makeSessionKey(sid)
	ret, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		hlog.Infof("GetSession(): key is nil, sid: %d", sid)
		return 0, nil
	} else if err != nil {
		hlog.Errorf("GetSession(): redis query error, err: %+v", err)
		return 0, errRedis
	}
	uid, err := strconv.ParseInt(ret, 10, 64)
	if err != nil {
		hlog.Errorf("GetSession(): parse user id error, err: %+v", err)
		return 0, errRedis
	}
	return uid, nil
}

func SetSession(ctx context.Context, sid *string, uid int64) berrors.Berror {
	key := makeSessionKey(sid)
	_, err := redisClient.Set(ctx, key, uid, 0).Result()
	if err != nil {
		hlog.Errorf("SetSession(): redis set error, err: %+v", err)
		return errRedis
	}
	return nil
}

func DelSession(ctx context.Context, sid *string) berrors.Berror {
	key := makeSessionKey(sid)
	n, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		hlog.Errorf("DelSession(): redis del error, err: %+v", err)
		return errRedis
	}
	if n == 0 {
		hlog.Errorf("DelSession(): no sid found, sid: %s", sid)
		return errRedis
	}
	return nil
}
