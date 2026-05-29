package rdb

import (
	"context"
	"errors"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UserJoinTeam(ctx context.Context, uid, tid int64) (bool, berrors.Berror) {
	trans := db.Begin()
	if trans.Error != nil {
		hlog.Errorf("UserJoinTeam(): begin trans failed, err: %+v", trans.Error)
		return false, errDB
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
			hlog.Infof("UserJoinTeam(): no user record, uid: %d", uid)
			return false, errNoUserRecord
		}
		hlog.Errorf("UserJoinTeam(): query user failed, err: %+v", err)
		return false, errDB
	}
	if usr.BelongTeam != 0 {
		trans.WithContext(ctx).Rollback()
		return false, errAlreadyInTeam
	}

	result := trans.WithContext(ctx).
		Table(tableTeamInfo).
		Where("team_id = ? AND member_cnt < ?", tid, contest.MaxTeamMember).
		UpdateColumn("member_cnt", gorm.Expr("member_cnt + 1"))
	if result.Error != nil {
		trans.WithContext(ctx).Rollback()
		hlog.Errorf("UserJoinTeam(): update member cnt trans failed, err: %+v", result.Error)
		return false, errDB
	}
	if result.RowsAffected == 0 {
		team := model.Team{}
		err = trans.WithContext(ctx).Table(tableTeamInfo).Where("team_id = ?", tid).Take(&team).Error
		trans.WithContext(ctx).Rollback()
		if err == nil {
			return false, nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hlog.Infof("UserJoinTeam(): no team record, tid: %d", tid)
			return false, errNoTeamRecord
		}
		hlog.Errorf("UserJoinTeam(): query team failed, err: %+v", err)
		return false, errDB
	}

	result = trans.WithContext(ctx).
		Table(tableUserInfo).
		Where("user_id = ? AND belong_team = 0", uid).
		Update("belong_team", tid)
	if result.Error != nil {
		hlog.Errorf("UserJoinTeam(): update user belong team trans failed, err: %+v", result.Error)
		trans.WithContext(ctx).Rollback()
		return false, errDB
	}
	if result.RowsAffected == 0 {
		hlog.Errorf("UserJoinTeam(): unexpected rows affected when updating user, uid: %d, tid: %d", uid, tid)
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
	trans := db.Begin()
	if trans.Error != nil {
		hlog.Errorf("UserQuitTeam(): begin trans failed, err: %+v", trans.Error)
		return errDB
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
			hlog.Infof("UserQuitTeam(): no user record, uid: %d", uid)
			return errNoUserRecord
		}
		hlog.Errorf("UserQuitTeam(): query user failed, err: %+v", err)
		return errDB
	}

	tid := usr.BelongTeam
	if tid == 0 {
		trans.WithContext(ctx).Rollback()
		hlog.Infof("UserQuitTeam(): user is not in a team, uid: %d", uid)
		return errNotInATeam
	}

	team := model.Team{}
	err = trans.WithContext(ctx).
		Table(tableTeamInfo).
		Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
		Where("team_id = ?", tid).
		Take(&team).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hlog.Errorf("UserQuitTeam(): no team record, tid: %d", tid)
			return errNoTeamRecord
		}
		hlog.Errorf("UserQuitTeam(): query team failed, err: %+v", err)
		return errDB
	}
	if team.MemberCnt < 1 {
		trans.WithContext(ctx).Rollback()
		hlog.Errorf("UserQuitTeam(): team member cnt < 1, tid: %d", tid)
		return errDB
	}

	result := trans.WithContext(ctx).
		Table(tableUserInfo).
		Where("user_id = ? AND belong_team = ?", uid, tid).
		Update("belong_team", 0)
	if result.Error != nil {
		hlog.Errorf("UserQuitTeam(): update user trans failed, err: %+v", result.Error)
		trans.WithContext(ctx).Rollback()
		return errDB
	}
	if result.RowsAffected == 0 {
		hlog.Errorf("UserQuitTeam(): unexpected rows affected when updating user, uid: %d, tid: %d", uid, tid)
		trans.WithContext(ctx).Rollback()
		return errDB
	}

	var resultTeam *gorm.DB
	if team.MemberCnt == 1 {
		resultTeam = trans.WithContext(ctx).Table(tableTeamInfo).Where("team_id = ?", tid).Delete(&model.Team{})
	} else {
		resultTeam = trans.WithContext(ctx).
			Table(tableTeamInfo).
			Where("team_id = ? AND member_cnt > 0", tid).
			UpdateColumn("member_cnt", gorm.Expr("member_cnt - 1"))
	}
	if resultTeam.Error != nil {
		hlog.Errorf("UserQuitTeam(): update team trans failed, err: %+v", resultTeam.Error)
		trans.WithContext(ctx).Rollback()
		return errDB
	}
	if resultTeam.RowsAffected != 1 {
		hlog.Errorf("UserQuitTeam(): unexpected rows affected when updating team, tid: %d, rows: %d", tid, resultTeam.RowsAffected)
		trans.WithContext(ctx).Rollback()
		return errDB
	}

	if err := trans.Commit().Error; err != nil {
		hlog.Errorf("UserQuitTeam(): transaction failed, err: %+v", err)
		return errDB
	}
	return nil
}
