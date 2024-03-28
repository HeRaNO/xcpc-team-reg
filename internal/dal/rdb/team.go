package rdb

import (
	"context"
	"errors"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

func GetUserInfosByTeamID(ctx context.Context, tid int64) ([]model.UserInfo, berrors.Berror) {
	rec := make([]model.User, contest.MaxTeamMember)
	result := db.Model(&model.User{}).Table(tableUserInfo).Where("belong_team = ?", tid).Find(&rec)

	if result.Error != nil {
		hlog.Errorf("GetUserInfosByTeamID(): query failed, err: %+v", result.Error)
		return nil, errDB
	}
	if result.RowsAffected == 0 {
		hlog.Errorf("GetUserInfosByTeamID(): no user in team, tid: %d", tid)
		return nil, errDB
	}
	if result.RowsAffected > int64(contest.MaxTeamMember) {
		hlog.Errorf("GetUserInfosByTeamID(): %d user in team, tid: %d", result.RowsAffected, tid)
		return nil, errDB
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

func CreateNewTeam(ctx context.Context, uid int64, team *model.TeamInfoModifyReq) (int64, string, berrors.Berror) {
	inviteToken, errTk := utils.GenToken(contest.UserTokenLength)
	if errTk != nil {
		hlog.Errorf("CreateNewTeam(): gen token failed: %+v", errTk.Msg())
		return 0, "", errTk
	}

	trans := db.Begin()
	info := model.Team{
		TeamName:        *team.TeamName,
		TeamAffiliation: *team.TeamAffiliation,
		MemberCnt:       1,
		InviteToken:     inviteToken,
	}
	err := trans.WithContext(ctx).Model(&model.Team{}).Table(tableTeamInfo).Create(&info).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, "", errDuplicatedName
		}
		hlog.Errorf("CreateNewTeam(): trans create team failed, err: %+v", err)
		return 0, "", errDB
	}

	tid := info.TeamID

	err = trans.WithContext(ctx).Table(tableUserInfo).Where("user_id = ?", uid).Update("belong_team", tid).Error
	if err != nil {
		hlog.Errorf("CreateNewTeam(): trans update user team id failed, err: %+v", err)
		trans.WithContext(ctx).Rollback()
		return 0, "", errDB
	}

	if err := trans.Commit().Error; err != nil {
		hlog.Errorf("CreateNewTeam(): transaction failed, err: %+v", err)
		return 0, "", errDB
	}

	return tid, inviteToken, nil
}

func GetTeamInfoByTeamID(ctx context.Context, tid int64) (*model.TeamInfo, berrors.Berror) {
	rec := make([]model.Team, 0)
	result := db.Model(&model.Team{}).Table(tableTeamInfo).Where("team_id = ?", tid).Find(&rec)

	if result.Error != nil {
		hlog.Errorf("GetTeamInfoByTeamID(): query failed, err: %+v", result.Error)
		return nil, errDB
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetTeamInfoByTeamID(): no team record, tid: %d", tid)
		return nil, errDB
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

func GetTeamInviteTokenByTeamID(ctx context.Context, tid int64) (*string, berrors.Berror) {
	ori := map[string]interface{}{}
	result := db.WithContext(ctx).Table(tableTeamInfo).Select("invite_token").Where("team_id = ?", tid).Find(&ori)

	if result.Error != nil {
		hlog.Errorf("GetTeamInviteTokenByTeamID(): query failed, err: %s", result.Error)
		return nil, errDB
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetTeamInviteTokenByTeamID(): no team record, tid: %d", tid)
		return nil, errNoTeamRecord
	}

	token := ori["invite_token"].(string)

	return &token, nil
}

func ModifyTeamInfoByTeamID(ctx context.Context, tid int64, team *model.TeamInfoModifyReq) berrors.Berror {
	result := db.WithContext(ctx).Model(&model.TeamInfoModifyReq{}).Table(tableTeamInfo).Where("team_id = ?", tid).Updates(team)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errDuplicatedName
		}
		hlog.Errorf("ModifyTeamInfoByTeamID(): update failed: %+v", result.Error)
		return errDB
	}
	if result.RowsAffected == 0 {
		hlog.Infof("ModifyTeamInfoByTeamID(): no record affected, tid: %d, info: %+v", tid, team)
	}

	return nil
}

func GetAllTeams() ([]model.Team, error) {
	teams := make([]model.Team, 0)
	err := db.Model(&model.Team{}).Table(tableTeamInfo).Find(&teams).Error
	if err != nil {
		return nil, err
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
