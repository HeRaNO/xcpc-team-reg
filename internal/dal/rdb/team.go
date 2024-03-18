package rdb

import (
	"context"
	"errors"

	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetUserInfosByTeamID(ctx context.Context, tid int64) ([]model.UserInfo, error) {
	rec := make([]model.User, contest.MaxTeamMember)
	result := db.Model(&model.User{}).Table(tableUserInfo).Where("belong_team = ?", tid).Find(&rec)

	if result.Error != nil {
		hlog.Errorf("GetUserInfosByTeamID(): query failed, err: %+v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		hlog.Errorf("GetUserInfosByTeamID(): no user in team, tid: %d", tid)
		return nil, errors.New("no user in this team but why???")
	}
	if result.RowsAffected > int64(contest.MaxTeamMember) {
		hlog.Errorf("GetUserInfosByTeamID(): %d user in team, tid: %d", result.RowsAffected, tid)
		return nil, errors.New("too many members in this team but why???")
	}

	usrInfo := make([]model.UserInfo, 0)

	for _, usr := range rec {
		usrInfo = append(usrInfo, model.UserInfo{
			Name:       usr.Name,
			School:     usr.School,
			StuID:      usr.StuID,
			BelongTeam: usr.BelongTeam,
			IsUESTCStu: usr.IsUESTCStu,
			Tshirt:     usr.Tshirt,
		})
	}

	return usrInfo, nil
}

func CreateNewTeam(ctx context.Context, uid int64, team *model.TeamInfoModifyReq) (int64, string, error) {
	err := redis.AddTeamName(ctx, team.TeamName)
	if err != nil {
		return 0, "", err
	}

	inviteToken, err := utils.GenToken(contest.UserTokenLength)
	if err != nil {
		return 0, "", err
	}

	trans := db.Begin()
	info := model.Team{
		TeamName:        *team.TeamName,
		TeamAffiliation: *team.TeamAffiliation,
		MemberCnt:       1,
		InviteToken:     inviteToken,
	}
	err = trans.WithContext(ctx).Model(&model.Team{}).Table(tableTeamInfo).Create(&info).Error
	if err != nil {
		hlog.Errorf("CreateNewTeam(): trans create team failed, err: %+v", err)
		redis.DelTeamName(ctx, team.TeamName)
		trans.WithContext(ctx).Rollback()
		return 0, "", err
	}

	tid := info.TeamID

	err = trans.WithContext(ctx).Table(tableUserInfo).Where("user_id = ?", uid).Update("belong_team", tid).Error
	if err != nil {
		hlog.Errorf("CreateNewTeam(): trans update user team id failed, err: %+v", err)
		trans.WithContext(ctx).Rollback()
		redis.DelTeamName(ctx, team.TeamName)
		return 0, "", err
	}

	if err := trans.Commit().Error; err != nil {
		hlog.Errorf("CreateNewTeam(): transaction failed, err: %+v", err)
		redis.DelTeamName(ctx, team.TeamName)
		return 0, "", err
	}

	return tid, inviteToken, nil
}

func GetTeamInfoByTeamID(ctx context.Context, tid int64) (*model.TeamInfo, error) {
	rec := make([]model.Team, 0)
	result := db.Model(&model.Team{}).Table(tableTeamInfo).Where("team_id = ?", tid).Find(&rec)

	if result.Error != nil {
		hlog.Errorf("GetTeamInfoByTeamID(): query failed, err: %+v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetTeamInfoByTeamID(): no team record, tid: %d", tid)
		return nil, errors.New("no team record")
	}

	usrInfo, err := GetUserInfosByTeamID(ctx, tid)
	if err != nil {
		return nil, err
	}

	teamInfo := &model.TeamInfo{
		TeamID:          tid,
		TeamName:        rec[0].TeamName,
		TeamAccount:     rec[0].TeamAccount,
		TeamPassword:    rec[0].TeamPassword,
		InviteToken:     rec[0].InviteToken,
		TeamMember:      usrInfo,
		MemberCnt:       rec[0].MemberCnt,
		TeamAffiliation: rec[0].TeamAffiliation,
	}
	return teamInfo, nil
}

func GetTeamInviteTokenByTeamID(ctx context.Context, tid int64) (*string, error) {
	ori := map[string]interface{}{}
	result := db.WithContext(ctx).Table(tableTeamInfo).Select("invite_token").Where("team_id = ?", tid).Find(&ori)

	if result.Error != nil {
		hlog.Errorf("GetTeamInviteTokenByTeamID(): query failed, err: %s", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetTeamInviteTokenByTeamID(): no team record, tid: %d", tid)
		return nil, errors.New("no team record")
	}

	token := ori["invite_token"].(string)

	return &token, nil
}

func GetTeamNameByTeamID(ctx context.Context, tid int64) (*string, error) {
	ori := map[string]interface{}{}
	result := db.WithContext(ctx).Table(tableTeamInfo).Select("team_name").Where("team_id = ?", tid).Find(&ori)

	if result.Error != nil {
		hlog.Errorf("GetTeamNameByTeamID(): query ori team name failed, err: %+v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetTeamNameByTeamID(): no team record, tid: %d", tid)
		return nil, errors.New("no team record")
	}

	oriTeamName := ori["team_name"].(string)
	return &oriTeamName, nil
}

func ModifyTeamInfoByTeamID(ctx context.Context, tid int64, team *model.TeamInfoModifyReq) error {
	trans := db.Begin()
	err := trans.WithContext(ctx).Model(&model.TeamInfoModifyReq{}).Table(tableTeamInfo).Where("team_id = ?", tid).Updates(team).Error
	if err != nil {
		hlog.Errorf("ModifyTeamInfoByTeamID(): trans make update failed: %+v", err)
		trans.WithContext(ctx).Rollback()
		return err
	}

	oriTeamName := new(string)

	if team.TeamName != nil {
		oriTeamName, err = GetTeamNameByTeamID(ctx, tid)
		if err != nil {
			trans.WithContext(ctx).Rollback()
			return err
		}
		err = redis.DelTeamName(ctx, oriTeamName)
		if err != nil {
			trans.WithContext(ctx).Rollback()
			return err
		}
		err = redis.AddTeamName(ctx, team.TeamName)
		if err != nil {
			trans.WithContext(ctx).Rollback()
			redis.AddTeamName(ctx, oriTeamName)
			return err
		}
	}

	if err := trans.Commit().Error; err != nil {
		hlog.Errorf("ModifyTeamInfoByTeamID(): transaction failed, err: %+v", err)
		if team.TeamName != nil {
			redis.DelTeamName(ctx, team.TeamName)
			redis.AddTeamName(ctx, oriTeamName)
		}
		return err
	}
	return nil
}

func GetAllTeams() ([]model.Team, error) {
	teams := make([]model.Team, 0)
	result := db.Model(&model.Team{}).Table(tableTeamInfo).Find(&teams)
	if result.Error != nil {
		return nil, result.Error
	}
	return teams, nil
}

func SetTeamAccPwdByID(tid int64, acc *string, pwd *string) error {
	accPwd := &model.Team{
		TeamAccount:  *acc,
		TeamPassword: *pwd,
	}
	return db.Model(&model.Team{}).Table(tableTeamInfo).Where("team_id = ?", tid).Updates(accPwd).Error
}
