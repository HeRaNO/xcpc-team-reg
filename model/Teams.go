package model

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/util"
	"github.com/go-redis/redis/v8"
)

const (
	TableTeamInfo  = "t_team"
	TeamNameSetKey = "TEAMNAME:SET"
)

type TeamInfoModify struct {
	TeamName string `json:"team_name"`
}

type TeamInfo struct {
	TeamName     string     `json:"teamname"`
	TeamAccount  string     `json:"account"`
	TeamPassword string     `json:"password"`
	InviteToken  string     `json:"invite_token"`
	TeamMember   []UserInfo `json:"member"`
}

type Team struct {
	TeamID       int64  `gorm:"column:team_id;primaryKey" json:"teamid"`
	TeamName     string `gorm:"column:team_name" json:"teamname"`
	MemberCnt    int    `gorm:"column:member_cnt" json:"memcnt"`
	TeamAccount  string `gorm:"column:team_account" json:"account"`
	TeamPassword string `gorm:"column:team_password" json:"password"`
}

func MakeTeamInviteTokenKey(tid int64) string {
	ret := fmt.Sprintf("TEAMINVITETOKEN:%d", tid)
	return ret
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

func GetTeamInviteToken(ctx context.Context, tid int64) (string, error) {
	key := MakeTeamInviteTokenKey(tid)
	ret, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[INFO] GetTeamInviteToken(): key is nil, tid: %d\n", tid)
		return "", nil
	} else if err != nil {
		log.Println("[ERROR] GetTeamInviteToken(): redis query error")
		return "", err
	}
	return ret, nil
}

func SetTeamInviteToken(ctx context.Context, tid int64, token *string) error {
	key := MakeTeamInviteTokenKey(tid)
	err := config.RedisClient.Set(ctx, key, *token, 0).Err()
	if err != nil {
		log.Println("[ERROR] SetTeamInviteToken(): redis set error")
		return err
	}
	return nil
}

func DelTeamInviteToken(ctx context.Context, tid int64) error {
	key := MakeTeamInviteTokenKey(tid)
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		log.Println("[ERROR] DelTeamInviteToken(): redis del error")
		return err
	}
	return nil
}

func CreateNewTeam(ctx context.Context, tname *string, uid int64) (string, error) {
	err := AddTeamName(ctx, tname)
	if err != nil {
		return "", err
	}

	trans := config.RDB.Begin()
	info := Team{
		TeamName:  *tname,
		MemberCnt: 1,
	}
	err = trans.WithContext(ctx).Model(&Team{}).Table(TableTeamInfo).Create(&info).Error
	if err != nil {
		DelTeamName(ctx, tname)
		trans.WithContext(ctx).Rollback()
		return "", err
	}

	tid := info.TeamID
	inviteToken := util.GenToken(config.UserTokenLength)
	err = SetTeamInviteToken(ctx, tid, &inviteToken)
	if err != nil {
		trans.WithContext(ctx).Rollback()
		DelTeamName(ctx, tname)
		return "", err
	}

	err = UpdateTeamIDByUserID(ctx, uid, tid)
	if err != nil {
		trans.WithContext(ctx).Rollback()
		DelTeamName(ctx, tname)
		return "", err
	}

	if err := trans.Commit().Error; err != nil {
		log.Println("[ERROR] CreateNewTeam(): transaction failed")
		DelTeamInviteToken(ctx, tid)
		DelTeamName(ctx, tname)
		return "", err
	}

	return inviteToken, nil
}

func GetTeamInfoByTeamID(ctx context.Context, tid int64) (*TeamInfo, error) {
	rdb := config.RDB

	rec := make([]Team, 0)
	result := rdb.Model(&Team{}).Table(TableTeamInfo).Where("team_id = ?", tid).Find(&rec)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no team record")
	}

	if result.RowsAffected > 1 {
		return nil, errors.New("duplicate team_id but why???")
	}

	inviteToken, err := GetTeamInviteToken(ctx, tid)
	if err != nil {
		return nil, err
	}

	teamInfo := new(TeamInfo)
	teamInfo.TeamName = rec[0].TeamName
	teamInfo.TeamAccount = rec[0].TeamAccount
	teamInfo.TeamPassword = rec[0].TeamPassword
	teamInfo.InviteToken = inviteToken

	usrInfo, err := GetUserInfosByTeamID(ctx, tid)
	if err != nil {
		return nil, err
	}

	teamInfo.TeamMember = usrInfo
	return teamInfo, nil
}

func ModifyTeamInfoByTeamID(ctx context.Context, tid int64, tname *string) error {
	trans := config.RDB.Begin()
	ori := map[string]interface{}{}
	result := trans.WithContext(ctx).Table(TableTeamInfo).Select("team_name").Where("team_id = ?", tid).Find(&ori)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no team record")
	}
	if result.RowsAffected > 1 {
		return errors.New("duplicate team_id but why???")
	}

	err := trans.WithContext(ctx).Table(TableTeamInfo).Where("team_id = ?", tid).Update("team_name", *tname).Error

	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}

	oriTeamName := ori["team_name"].(string)

	err = DelTeamName(ctx, &oriTeamName)
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}

	err = AddTeamName(ctx, tname)
	if err != nil {
		trans.WithContext(ctx).Rollback()
		AddTeamName(ctx, &oriTeamName)
		return err
	}

	if err := trans.Commit().Error; err != nil {
		log.Println("[ERROR] ModifyTeamInfoByTeamID(): transaction failed")
		DelTeamName(ctx, tname)
		AddTeamName(ctx, &oriTeamName)
		return err
	}
	return nil
}
