package redis

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

func IsTeamNameUsed(ctx context.Context, tname *string) (bool, error) {
	redisLock := newRedisMutex(redisClient, teamNameSetKey, 1, 200, 5)
	if !redisLock.lock(ctx) {
		hlog.Errorf("IsTeamNameUsed(): get mutex lock failed")
		return true, errors.New("get mutex lock failed")
	}
	defer redisLock.unlock(ctx)
	ret, err := redisClient.SIsMember(ctx, teamNameSetKey, *tname).Result()
	if err == redis.Nil {
		hlog.Info("IsTeamNameUsed(): key is nil")
		return false, nil
	} else if err != nil {
		hlog.Errorf("IsTeamNameUsed(): redis query error, err: %+v", err)
		return true, err
	}
	return ret, nil
}

func AddTeamName(ctx context.Context, tname *string) error {
	redisLock := newRedisMutex(redisClient, teamNameSetKey, 1, 200, 5)
	if !redisLock.lock(ctx) {
		hlog.Errorf("AddTeamName(): get mutex lock failed")
		return errors.New("get mutex lock failed")
	}
	defer redisLock.unlock(ctx)
	_, err := redisClient.SAdd(ctx, teamNameSetKey, *tname).Result()
	if err != nil {
		hlog.Errorf("AddTeamName(): redis sadd error, err: %+v", err)
		return err
	}
	return nil
}

func DelTeamName(ctx context.Context, tname *string) error {
	redisLock := newRedisMutex(redisClient, teamNameSetKey, 1, 200, 5)
	if !redisLock.lock(ctx) {
		hlog.Errorf("DelTeamName(): get mutex lock failed")
		return errors.New("get mutex lock failed")
	}
	defer redisLock.unlock(ctx)
	_, err := redisClient.SRem(ctx, teamNameSetKey, *tname).Result()
	if err != nil {
		hlog.Errorf("DelTeamName(): redis srem error, err: %+v", err)
		return err
	}
	return nil
}
