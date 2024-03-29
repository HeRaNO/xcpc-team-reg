package redis

import (
	"context"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

func GetEmailToken(ctx context.Context, email *string) (string, berrors.Berror) {
	key := makeEmailTokenKey(email)
	ret, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		hlog.Infof("GetEmailToken(): key is nil, email: %s", *email)
		return "", nil
	} else if err != nil {
		hlog.Errorf("GetEmailToken(): redis query error, err: %+v", err)
		return "", errRedis
	}
	return ret, nil
}

func SetEmailToken(ctx context.Context, email *string, token *string, exptime time.Duration) berrors.Berror {
	key := makeEmailTokenKey(email)
	err := redisClient.Set(ctx, key, *token, exptime).Err()
	if err != nil {
		hlog.Errorf("SetEmailToken(): redis set error, err: %+v", err)
		return errRedis
	}
	return nil
}

func DelEmailToken(ctx context.Context, email *string) berrors.Berror {
	key := makeEmailTokenKey(email)
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		hlog.Errorf("DelEmailToken(): redis del error, err: %+v", err)
		return errRedis
	}
	return nil
}

func GetEmailRequest(ctx context.Context, email *string) berrors.Berror {
	key := makeEmailRequestKey(email)
	err := redisClient.Get(ctx, key).Err()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		hlog.Errorf("GetEmailRequest(): redis query error, err: %+v", err)
		return errRedis
	}
	return errTooFrequent
}

func SetEmailRequest(ctx context.Context, email *string, exptime time.Duration) berrors.Berror {
	key := makeEmailRequestKey(email)
	err := redisClient.Set(ctx, key, "1", exptime).Err()
	if err != nil {
		hlog.Errorf("SetEmailRequest(): redis set error, err: %+v", err)
		return errRedis
	}
	return nil
}

func GetEmailAction(ctx context.Context, email *string) (string, berrors.Berror) {
	key := makeEmailActionKey(email)
	ret, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		hlog.Errorf("GetEmailRequest(): redis query error, err: %+v", err)
		return "", errRedis
	}
	return ret, nil
}

func SetEmailAction(ctx context.Context, email *string, action *string, exptime time.Duration) berrors.Berror {
	key := makeEmailActionKey(email)
	err := redisClient.Set(ctx, key, *action, exptime).Err()
	if err != nil {
		hlog.Errorf("SetEmailAction(): redis set error, err: %+v", err)
		return errRedis
	}
	return nil
}

func DelEmailAction(ctx context.Context, email *string) berrors.Berror {
	key := makeEmailActionKey(email)
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		hlog.Errorf("DelEmailAction(): redis del error, err: %+v", err)
		return errRedis
	}
	return nil
}
