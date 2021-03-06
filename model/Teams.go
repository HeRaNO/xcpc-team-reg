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

type JoinTeamRequest struct {
	TeamID      string `json:"team_id"`
	InviteToken string `json:"invite_token"`
}

type TeamInfo struct {
	TeamID       int64      `json:"team_id"`
	TeamName     string     `json:"team_name"`
	TeamAccount  string     `json:"account"`
	TeamPassword string     `json:"password"`
	InviteToken  string     `json:"invite_token"`
	TeamMember   []UserInfo `json:"member"`
	MemberCnt    int        `json:"mem_cnt"`
}

type Team struct {
	TeamID       int64  `gorm:"column:team_id;primaryKey" json:"team_id"`
	TeamName     string `gorm:"column:team_name" json:"team_name"`
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

func GetTeamCount(ctx context.Context) (int64, error) {
	redisLock := util.NewRedisMutex(config.RedisClient, TeamNameSetKey, 1, 200, 5)
	redisLock.Lock(ctx)
	if !redisLock.Lock(ctx) {
		log.Printf("[ERROR] GetTeamCount(): get mutex lock failed\n")
		return 0, errors.New("get mutex lock failed")
	}
	defer redisLock.Unlock(ctx)
	cnt, err := config.RedisClient.SCard(ctx, TeamNameSetKey).Result()
	if err != nil {
		log.Println("[ERROR] GetTeamCount(): redis scard error")
		return 0, err
	}
	return cnt, nil
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

func ValidateTeamInviteToken(ctx context.Context, tid int64, token *string) (bool, error) {
	tokenFromRedis, err := GetTeamInviteToken(ctx, tid)
	if tokenFromRedis == "" || err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if *token != tokenFromRedis {
		return false, nil
	}
	return true, nil
}

func CreateNewTeam(ctx context.Context, tname *string, uid int64) (int64, string, error) {
	err := AddTeamName(ctx, tname)
	if err != nil {
		return -1, "", err
	}

	inviteToken, err := util.GenToken(config.UserTokenLength)
	if err != nil {
		return -1, "", err
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
		return -1, "", err
	}

	tid := info.TeamID
	err = SetTeamInviteToken(ctx, tid, &inviteToken)
	if err != nil {
		trans.WithContext(ctx).Rollback()
		DelTeamName(ctx, tname)
		return -1, "", err
	}

	err = trans.WithContext(ctx).Table(TableUserInfo).Where("user_id = ?", uid).Update("belong_team", tid).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		DelTeamName(ctx, tname)
		return -1, "", err
	}

	if err := trans.Commit().Error; err != nil {
		log.Println("[ERROR] CreateNewTeam(): transaction failed")
		DelTeamInviteToken(ctx, tid)
		DelTeamName(ctx, tname)
		return -1, "", err
	}

	return tid, inviteToken, nil
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
	teamInfo.TeamID = tid
	teamInfo.TeamName = rec[0].TeamName
	teamInfo.MemberCnt = rec[0].MemberCnt
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

func UserJoinTeam(ctx context.Context, uid int64, tid int64) (bool, error) {
	trans := config.RDB.Begin()
	ori := map[string]interface{}{}
	result := trans.WithContext(ctx).Table(TableTeamInfo).Select("member_cnt").Where("team_id = ?", tid).Find(&ori)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, errors.New("no team record")
	}
	if result.RowsAffected > 1 {
		return false, errors.New("duplicate team_id but why???")
	}

	nowCnt := ori["member_cnt"].(int32)
	if nowCnt >= config.MaxTeamMember {
		return false, nil
	}
	if nowCnt < 1 {
		return false, errors.New("the member cnt is less than 1 but why???")
	}

	nowCnt++

	err := trans.WithContext(ctx).Table(TableTeamInfo).Where("team_id = ?", tid).Update("member_cnt", nowCnt).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return false, err
	}

	err = trans.WithContext(ctx).Table(TableUserInfo).Where("user_id = ?", uid).Update("belong_team", tid).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return false, err
	}

	if err := trans.Commit().Error; err != nil {
		log.Println("[ERROR] UserJoinTeam(): transaction failed")
		return false, err
	}
	return true, nil
}

func UserQuitTeam(ctx context.Context, uid int64, tid int64) error {
	trans := config.RDB.Begin()
	ori := make([]Team, 0)
	result := trans.WithContext(ctx).Table(TableTeamInfo).Where("team_id = ?", tid).Find(&ori)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no team record")
	}
	if result.RowsAffected > 1 {
		return errors.New("duplicate team_id but why???")
	}

	team := ori[0]
	nowCnt := team.MemberCnt

	nowCnt--

	err := trans.WithContext(ctx).Table(TableTeamInfo).Where("team_id = ?", tid).Update("member_cnt", nowCnt).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}

	if nowCnt < 1 {
		trans.WithContext(ctx).Table(TableTeamInfo).Where("team_id = ?", tid).Delete(&Team{})
		if err != nil {
			trans.WithContext(ctx).Rollback()
			return err
		}
	}

	err = trans.WithContext(ctx).Table(TableUserInfo).Where("user_id = ?", uid).Update("belong_team", 0).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}

	inviteToken := ""
	if nowCnt < 1 {
		err := DelTeamName(ctx, &team.TeamName)
		if err != nil {
			trans.WithContext(ctx).Rollback()
			return err
		}
		inviteToken, err = GetTeamInviteToken(ctx, tid)
		if err != nil {
			trans.WithContext(ctx).Rollback()
			return err
		}
		err = DelTeamInviteToken(ctx, tid)
		if err != nil {
			trans.WithContext(ctx).Rollback()
			AddTeamName(ctx, &team.TeamName)
			return err
		}
	}

	if err := trans.Commit().Error; err != nil {
		log.Println("[ERROR] UserJoinTeam(): transaction failed")
		if nowCnt < 1 {
			AddTeamName(ctx, &team.TeamName)
			SetTeamInviteToken(ctx, tid, &inviteToken)
		}
		return err
	}
	return nil
}

func GetAllTeamIDs(ctx context.Context) ([]int64, error) {
	rdb := config.RDB

	cnt, _ := GetTeamCount(ctx)
	rec := make([]int64, cnt)
	result := rdb.WithContext(ctx).Table(TableTeamInfo).Select("team_id").Scan(&rec)

	if result.Error != nil {
		return nil, result.Error
	}

	return rec, nil
}

func SetTeamAccPwdByID(ctx context.Context, tid int64, acc *string, pwd *string) error {
	trans := config.RDB.Begin()
	accPwd := &Team{
		TeamAccount:  *acc,
		TeamPassword: *pwd,
	}
	result := trans.WithContext(ctx).Table(TableTeamInfo).Where("team_id = ?", tid).Updates(accPwd)
	row, err := result.RowsAffected, result.Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}
	if row < 1 {
		trans.WithContext(ctx).Rollback()
		return errors.New("no team record found")
	}
	if row > 1 {
		trans.WithContext(ctx).Rollback()
		return errors.New("duplicate team_id but why???")
	}
	if err := trans.Commit().Error; err != nil {
		log.Println("[ERROR] SetTeamAccPwdByID: transaction failed")
		return err
	}
	return nil
}
