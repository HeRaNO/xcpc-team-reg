package rdb

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func UserJoinTeam(ctx context.Context, uid int64, tid int64) (bool, berrors.Berror) {
	ori := map[string]interface{}{}
	result := db.WithContext(ctx).Table(tableTeamInfo).Select("member_cnt").Where("team_id = ?", tid).Find(&ori)

	if result.Error != nil {
		hlog.Errorf("UserJoinTeam(): query member cnt failed, err: %+v", result.Error)
		return false, errDB
	}
	if result.RowsAffected == 0 {
		hlog.Infof("UserJoinTeam(): no team record, tid: %d", tid)
		return false, errNoTeamRecord
	}

	nowCnt := ori["member_cnt"].(int32)
	if nowCnt == contest.MaxTeamMember {
		return false, nil
	}
	if nowCnt > contest.MaxTeamMember {
		hlog.Errorf("UserJoinTeam(): team member cnt > %d, tid: %d", contest.MaxTeamMember, tid)
		return false, errDB
	}
	if nowCnt < 1 {
		hlog.Errorf("UserJoinTeam(): team member cnt < 1, tid: %d", tid)
		return false, errDB
	}

	nowCnt++
	trans := db.Begin()
	err := trans.WithContext(ctx).Table(tableTeamInfo).Where("team_id = ?", tid).Update("member_cnt", nowCnt).Error
	if err != nil {
		hlog.Errorf("UserJoinTeam(): update member cnt trans failed, err: %+v", err)
		trans.WithContext(ctx).Rollback()
		return false, errDB
	}

	err = trans.WithContext(ctx).Table(tableUserInfo).Where("user_id = ?", uid).Update("belong_team", tid).Error
	if err != nil {
		hlog.Errorf("UserJoinTeam(): update user belong team trans failed, err: %+v", err)
		trans.WithContext(ctx).Rollback()
		return false, errDB
	}

	if err := trans.Commit().Error; err != nil {
		hlog.Errorf("UserJoinTeam(): transaction failed, err: %+v", err)
		return false, errDB
	}
	return true, nil
}

func UserQuitTeam(ctx context.Context, uid int64) berrors.Berror {
	usrinfo, err := GetUserInfoByID(ctx, uid)
	if err != nil {
		hlog.Errorf("UserQuitTeam(): query user failed, err: %+v", err)
		return err
	}

	tid := usrinfo.BelongTeam
	if tid == 0 {
		hlog.Infof("UserQuitTeam(): user is not in a team, uid: %d", uid)
		return errNotInATeam
	}

	ori := make([]model.Team, 0)
	result := db.WithContext(ctx).Table(tableTeamInfo).Where("team_id = ?", tid).Find(&ori)

	if result.Error != nil {
		hlog.Errorf("UserQuitTeam(): query team failed, err: %+v", result.Error)
		return errDB
	}
	if result.RowsAffected == 0 {
		hlog.Errorf("UserQuitTeam(): no team record, tid: %d", tid)
		return errNoTeamRecord
	}

	team := ori[0]
	nowCnt := team.MemberCnt - 1

	trans := db.Begin()

	var errDel error
	if nowCnt < 1 {
		errDel = trans.WithContext(ctx).Table(tableTeamInfo).Where("team_id = ?", tid).Delete(&model.Team{}).Error
	} else {
		errDel = trans.WithContext(ctx).Table(tableTeamInfo).Where("team_id = ?", tid).Update("member_cnt", nowCnt).Error
	}
	if errDel != nil {
		hlog.Errorf("UserQuitTeam(): update team trans failed, err: %+v", errDel)
		trans.WithContext(ctx).Rollback()
		return errDB
	}

	errBel := trans.WithContext(ctx).Table(tableUserInfo).Where("user_id = ?", uid).Update("belong_team", 0).Error
	if errBel != nil {
		hlog.Errorf("UserQuitTeam(): update user trans failed, err: %+v", err)
		trans.WithContext(ctx).Rollback()
		return errDB
	}

	if err := trans.Commit().Error; err != nil {
		hlog.Errorf("UserQuitTeam(): transaction failed, err: %+v", err)
		return errDB
	}
	return nil
}
