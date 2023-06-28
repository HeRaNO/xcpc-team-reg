package redis

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

func GetUserIDByEmail(ctx context.Context, email *string) (int64, error) {
	key := makeEmailUserIDKey(email)
	ret, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		hlog.Infof("GetUserIDByEmail(): key is nil, email: %s", *email)
		return 0, nil
	} else if err != nil {
		hlog.Errorf("GetUserIDByEmail(): redis query error, err: %+v", err)
		return 0, err
	}
	uid, err := strconv.ParseInt(ret, 10, 64)
	if err != nil {
		hlog.Errorf("GetUserIDByEmail(): parse uid failed, err: %+v", err)
		return 0, err
	}
	return uid, nil
}

func SetEmailUserID(ctx context.Context, uid int64, email *string) error {
	key := makeEmailUserIDKey(email)
	err := redisClient.Set(ctx, key, uid, 0).Err()
	if err != nil {
		hlog.Errorf("SetEmailUserID(): redis set error, err %+v", err)
		return err
	}
	return nil
}

func DelEmailUserID(ctx context.Context, email *string) error {
	key := makeEmailUserIDKey(email)
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		hlog.Errorf("DelEmailUserID(): redis del error, err: %+v", err)
		return err
	}
	return nil
}
