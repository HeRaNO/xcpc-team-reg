package model

import (
	"context"
	"errors"
	"log"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/util"
	"github.com/go-redis/redis/v8"
)

const (
	TableTeamInfo  = "t_team"
	TeamNameSetKey = "TEAMNAME:SET"
)

type Team struct {
	TeamID       int64  `gorm:"column:team_id primaryKey" json:"teamid"`
	TeamName     string `gorm:"column:team_name" json:"teamname"`
	MemberCnt    int    `gorm:"column:member_cnt" json:"memcnt"`
	TeamAccount  string `gorm:"column:team_account" json:"account"`
	TeamPassword string `gorm:"column:team_password" json:"password"`
}

func IsTeamNameUsed(ctx context.Context, tname *string) (bool, error) {
	redisLock := util.NewRedisMutex(config.RedisClient, TeamNameSetKey, 1, 200, 5)
	if !redisLock.Lock(ctx) {
		log.Printf("[ERROR] IsTeamNameUsed(): get mutex lock failed\n")
		return true, errors.New("get mutex lock failed")
	}
	defer redisLock.Unlock(ctx)
	ret, err := config.RedisClient.SIsMember(ctx, TeamNameSetKey, *tname).Result()
	if err == redis.Nil {
		log.Printf("[INFO] IsTeamNameUsed(): key is nil\n")
		return false, nil
	} else if err != nil {
		log.Println("[ERROR] IsTeamNameUsed(): redis query error")
		return true, err
	}
	return ret, nil
}

func AddTeamName(ctx context.Context, tname *string) error {
	redisLock := util.NewRedisMutex(config.RedisClient, TeamNameSetKey, 1, 200, 5)
	redisLock.Lock(ctx)
	if !redisLock.Lock(ctx) {
		log.Printf("[ERROR] AddTeamName(): get mutex lock failed\n")
		return errors.New("get mutex lock failed")
	}
	defer redisLock.Unlock(ctx)
	_, err := config.RedisClient.SAdd(ctx, TeamNameSetKey, *tname).Result()
	if err != nil {
		log.Println("[ERROR] AddTeamName(): redis sadd error")
		return err
	}
	return nil
}

func DelTeamName(ctx context.Context, tname *string) error {
	redisLock := util.NewRedisMutex(config.RedisClient, TeamNameSetKey, 1, 200, 5)
	redisLock.Lock(ctx)
	if !redisLock.Lock(ctx) {
		log.Printf("[ERROR] DelTeamName(): get mutex lock failed\n")
		return errors.New("get mutex lock failed")
	}
	defer redisLock.Unlock(ctx)
	_, err := config.RedisClient.SRem(ctx, TeamNameSetKey, *tname).Result()
	if err != nil {
		log.Println("[ERROR] DelTeamName(): redis srem error")
		return err
	}
	return nil
}
