package redis

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

func GetEmailToken(ctx context.Context, email *string) (string, error) {
	key := makeEmailTokenKey(email)
	ret, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		hlog.Infof("GetEmailToken(): key is nil, email: %s", *email)
		return "", nil
	} else if err != nil {
		hlog.Errorf("GetEmailToken(): redis query error, err: %+v", err)
		return "", err
	}
	return ret, nil
}

func SetEmailToken(ctx context.Context, email *string, token *string, exptime time.Duration) error {
	key := makeEmailTokenKey(email)
	err := redisClient.Set(ctx, key, *token, exptime).Err()
	if err != nil {
		hlog.Errorf("SetEmailToken(): redis set error, err: %+v", err)
		return err
	}
	return nil
}

func DelEmailToken(ctx context.Context, email *string) error {
	key := makeEmailTokenKey(email)
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		hlog.Errorf("DelEmailToken(): redis del error, err: %+v", err)
		return err
	}
	return nil
}

func GetEmailRequest(ctx context.Context, email *string) error {
	key := makeEmailRequestKey(email)
	err := redisClient.Get(ctx, key).Err()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		hlog.Errorf("GetEmailRequest(): redis query error, err: %+v", err)
		return err
	}
	return errors.New("email request too frequent")
}

func SetEmailRequest(ctx context.Context, email *string, exptime time.Duration) error {
	key := makeEmailRequestKey(email)
	err := redisClient.Set(ctx, key, "1", exptime).Err()
	if err != nil {
		hlog.Errorf("SetEmailRequest(): redis set error, err: %+v", err)
		return err
	}
	return nil
}

func GetEmailAction(ctx context.Context, email *string) (string, error) {
	key := makeEmailActionKey(email)
	ret, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		hlog.Errorf("GetEmailRequest(): redis query error, err: %+v", err)
		return "", err
	}
	return ret, nil
}

func SetEmailAction(ctx context.Context, email *string, action *string, exptime time.Duration) error {
	key := makeEmailActionKey(email)
	err := redisClient.Set(ctx, key, *action, exptime).Err()
	if err != nil {
		hlog.Errorf("SetEmailAction(): redis set error, err: %+v", err)
		return err
	}
	return nil
}

func DelEmailAction(ctx context.Context, email *string) error {
	key := makeEmailActionKey(email)
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		hlog.Errorf("DelEmailAction(): redis del error, err: %+v", err)
		return err
	}
	return nil
}
