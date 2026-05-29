package rdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserInfosByTeamID(ctx context.Context, tid int64) ([]model.UserInfo, berrors.Berror) {
	usrInfo, err := getUserInfosByTeamIDQuery(db.WithContext(ctx), tid)
	if err != nil {
		hlog.Errorf("GetUserInfosByTeamID(): %+v", err)
		return nil, errDB
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
	if trans.Error != nil {
		hlog.Errorf("CreateNewTeam(): begin trans failed, err: %+v", trans.Error)
		return 0, "", errDB
	}

	usr := model.User{}
	err := trans.WithContext(ctx).
		Table(tableUserInfo).
		Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
		Where("user_id = ?", uid).
		Take(&usr).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hlog.Infof("CreateNewTeam(): no user record, uid: %d", uid)
			return 0, "", errNoUserRecord
		}
		hlog.Errorf("CreateNewTeam(): query user failed, err: %+v", err)
		return 0, "", errDB
	}
	if usr.BelongTeam != 0 {
		trans.WithContext(ctx).Rollback()
		return 0, "", errAlreadyInTeam
	}

	info := model.Team{
		TeamName:        *team.TeamName,
		TeamAffiliation: *team.TeamAffiliation,
		MemberCnt:       1,
		InviteToken:     inviteToken,
	}
	err = trans.WithContext(ctx).Model(&model.Team{}).Table(tableTeamInfo).Create(&info).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, "", errDuplicatedName
		}
		hlog.Errorf("CreateNewTeam(): trans create team failed, err: %+v", err)
		return 0, "", errDB
	}

	tid := info.TeamID

	result := trans.WithContext(ctx).
		Table(tableUserInfo).
		Where("user_id = ? AND belong_team = 0", uid).
		Update("belong_team", tid)
	if result.Error != nil {
		hlog.Errorf("CreateNewTeam(): trans update user team id failed, err: %+v", result.Error)
		trans.WithContext(ctx).Rollback()
		return 0, "", errDB
	}
	if result.RowsAffected != 1 {
		hlog.Errorf("CreateNewTeam(): unexpected rows affected when updating user, uid: %d, rows: %d", uid, result.RowsAffected)
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
	trans, err := beginConsistentRead(ctx)
	if err != nil {
		hlog.Errorf("GetTeamInfoByTeamID(): begin read trans failed, err: %+v", err)
		return nil, errDB
	}

	rec := make([]model.Team, 0)
	result := trans.Table(tableTeamInfo).Where("team_id = ?", tid).Find(&rec)

	if result.Error != nil {
		trans.WithContext(ctx).Rollback()
		hlog.Errorf("GetTeamInfoByTeamID(): query failed, err: %+v", result.Error)
		return nil, errDB
	}
	if result.RowsAffected == 0 {
		trans.WithContext(ctx).Rollback()
		hlog.Infof("GetTeamInfoByTeamID(): no team record, tid: %d", tid)
		return nil, errNoTeamRecord
	}

	usrInfo, err := getUserInfosByTeamIDQuery(trans, tid)
	if err != nil {
		trans.WithContext(ctx).Rollback()
		hlog.Errorf("GetTeamInfoByTeamID(): %+v", err)
		return nil, errDB
	}
	if err := validateRecordedTeamMembers(tid, rec[0].MemberCnt, len(usrInfo)); err != nil {
		trans.WithContext(ctx).Rollback()
		hlog.Errorf("GetTeamInfoByTeamID(): %+v", err)
		return nil, errDB
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
	if err := trans.Commit().Error; err != nil {
		hlog.Errorf("GetTeamInfoByTeamID(): commit read trans failed, err: %+v", err)
		return nil, errDB
	}
	return teamInfo, nil
}

func GetTeamInviteTokenByTeamID(ctx context.Context, tid int64) (*string, berrors.Berror) {
	token := ""
	result := db.WithContext(ctx).Table(tableTeamInfo).Select("invite_token").Where("team_id = ?", tid).Scan(&token)

	if result.Error != nil {
		hlog.Errorf("GetTeamInviteTokenByTeamID(): query failed, err: %s", result.Error)
		return nil, errDB
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetTeamInviteTokenByTeamID(): no team record, tid: %d", tid)
		return nil, errNoTeamRecord
	}

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

func GetAllTeamInfos(ctx context.Context) ([]model.TeamInfo, error) {
	trans, err := beginConsistentRead(ctx)
	if err != nil {
		return nil, err
	}

	teams := make([]model.Team, 0)
	err = trans.Model(&model.Team{}).Table(tableTeamInfo).Order("team_id asc").Find(&teams).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return nil, err
	}
	if len(teams) == 0 {
		trans.WithContext(ctx).Rollback()
		return []model.TeamInfo{}, nil
	}

	teamIDs := make([]int64, 0, len(teams))
	for _, team := range teams {
		teamIDs = append(teamIDs, team.TeamID)
	}

	users := make([]model.User, 0)
	err = trans.
		Model(&model.User{}).
		Table(tableUserInfo).
		Where("belong_team IN ?", teamIDs).
		Order("belong_team asc, user_id asc").
		Find(&users).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return nil, err
	}

	memberMap := make(map[int64][]model.UserInfo, len(teams))
	for _, usr := range users {
		memberMap[usr.BelongTeam] = append(memberMap[usr.BelongTeam], buildUserInfo(usr))
	}

	teamInfos := make([]model.TeamInfo, 0, len(teams))
	for _, team := range teams {
		members := memberMap[team.TeamID]
		if err := validateRecordedTeamMembers(team.TeamID, team.MemberCnt, len(members)); err != nil {
			trans.WithContext(ctx).Rollback()
			hlog.Errorf("GetAllTeamInfos(): %+v", err)
			return nil, err
		}
		teamInfos = append(teamInfos, model.TeamInfo{
			TeamID:          team.TeamID,
			TeamName:        team.TeamName,
			TeamAccount:     team.TeamAccount,
			TeamPassword:    team.TeamPassword,
			InviteToken:     team.InviteToken,
			TeamMember:      members,
			MemberCnt:       team.MemberCnt,
			TeamAffiliation: team.TeamAffiliation,
		})
	}

	if err := trans.Commit().Error; err != nil {
		return nil, err
	}
	return teamInfos, nil
}

func beginConsistentRead(ctx context.Context) (*gorm.DB, error) {
	trans := db.WithContext(ctx).Begin(&sql.TxOptions{
		ReadOnly:  true,
		Isolation: sql.LevelRepeatableRead,
	})
	if trans.Error != nil {
		return nil, trans.Error
	}
	return trans, nil
}

func getUserInfosByTeamIDQuery(query *gorm.DB, tid int64) ([]model.UserInfo, error) {
	rec := make([]model.User, 0, contest.MaxTeamMember)
	result := query.Model(&model.User{}).Table(tableUserInfo).Where("belong_team = ?", tid).Find(&rec)
	if result.Error != nil {
		return nil, result.Error
	}
	if err := validateActualTeamMembers(tid, len(rec)); err != nil {
		return nil, err
	}

	usrInfo := make([]model.UserInfo, 0, len(rec))
	for _, usr := range rec {
		usrInfo = append(usrInfo, buildUserInfo(usr))
	}
	return usrInfo, nil
}

func buildUserInfo(usr model.User) model.UserInfo {
	return model.UserInfo{
		Name:       usr.Name,
		School:     usr.School,
		StuID:      usr.StuID,
		BelongTeam: usr.BelongTeam,
		IsUESTCStu: usr.IsUESTCStu,
		Tshirt:     usr.Tshirt,
	}
}

func validateActualTeamMembers(tid int64, actualCnt int) error {
	if actualCnt == 0 {
		return fmt.Errorf("team has no user record, tid: %d", tid)
	}
	if actualCnt > contest.MaxTeamMember {
		return fmt.Errorf("team member cnt exceeds limit, tid: %d, actual: %d", tid, actualCnt)
	}
	return nil
}

func validateRecordedTeamMembers(tid int64, recordedCnt, actualCnt int) error {
	if err := validateActualTeamMembers(tid, actualCnt); err != nil {
		return err
	}
	if recordedCnt != actualCnt {
		return fmt.Errorf("team member cnt mismatched, tid: %d, recorded: %d, actual: %d", tid, recordedCnt, actualCnt)
	}
	return nil
}

func SetTeamAccPwdByID(tid int64, acc, pwd *string) error {
	accPwd := &model.Team{
		TeamAccount:  *acc,
		TeamPassword: *pwd,
	}
	return db.Model(&model.Team{}).Table(tableTeamInfo).Where("team_id = ?", tid).Updates(accPwd).Error
}
