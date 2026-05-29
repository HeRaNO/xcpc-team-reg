package rdb

import (
	"context"
	"errors"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

func CreateNewUser(ctx context.Context, usr *UserRegInfo) berrors.Berror {
	trans := db.Begin()
	if trans.Error != nil {
		hlog.Errorf("CreateNewUser(): begin trans failed, err: %+v", trans.Error)
		return errDB
	}

	info := model.User{
		Name:       usr.Name,
		Email:      usr.Email,
		School:     usr.School,
		StuID:      usr.StuID,
		IsUESTCStu: usr.IsUESTCStu,
		Tshirt:     usr.Tshirt,
	}
	err := trans.WithContext(ctx).Model(&model.User{}).Table(tableUserInfo).Create(&info).Error
	if err != nil {
		hlog.Errorf("CreateNewUser(): create trans failed, err: %+v", err)
		trans.WithContext(ctx).Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errDuplicatedUser
		}
		return errDB
	}

	uid := info.UserID

	authInfo := &model.Auth{
		UserID: uid,
		Email:  usr.Email,
		Pwd:    usr.PwdToken,
	}
	err = trans.WithContext(ctx).Model(&model.Auth{}).Table(tableAuthInfo).Create(authInfo).Error
	if err != nil {
		hlog.Errorf("CreateNewUser(): create auth trans failed, err: %+v", err)
		trans.WithContext(ctx).Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errDuplicatedUser
		}
		return errDB
	}

	if err := trans.Commit().Error; nil != err {
		hlog.Errorf("CreateNewUser(): transaction failed, err: %+v", err)
		return errDB
	}

	if err := syncEmailUserIDCache(ctx, uid, &usr.Email); err != nil {
		hlog.Errorf("CreateNewUser(): sync email cache failed after commit, uid: %d, email: %s, err: %+v", uid, usr.Email, err)
	}

	return nil
}

func GetUserIDByEmail(ctx context.Context, email *string) (int64, berrors.Berror) {
	uid, err := redis.GetUserIDByEmail(ctx, email)
	if err == nil && uid != 0 {
		return uid, nil
	}
	if err != nil {
		hlog.Errorf("GetUserIDByEmail(): cache lookup failed, fallback to db, email: %s, err: %+v", *email, err)
	}

	usr := model.User{}
	dbErr := db.WithContext(ctx).
		Model(&model.User{}).
		Table(tableUserInfo).
		Select("user_id", "email").
		Where("email = ?", *email).
		Take(&usr).Error
	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		hlog.Errorf("GetUserIDByEmail(): db lookup failed, email: %s, err: %+v", *email, dbErr)
		return 0, errDB
	}

	if err := syncEmailUserIDCache(ctx, usr.UserID, email); err != nil {
		hlog.Errorf("GetUserIDByEmail(): sync email cache failed, uid: %d, email: %s, err: %+v", usr.UserID, *email, err)
	}

	return usr.UserID, nil
}

func syncEmailUserIDCache(ctx context.Context, uid int64, email *string) berrors.Berror {
	return redis.SetEmailUserID(ctx, uid, email)
}

func GetUserInfoByID(ctx context.Context, uid int64) (*model.UserInfo, berrors.Berror) {
	rec := make([]model.User, 0)
	result := db.Model(&model.User{}).Table(tableUserInfo).Where("user_id = ?", uid).Find(&rec)

	if result.Error != nil {
		hlog.Errorf("GetUserInfoByID(): query failed, err: %+v", result.Error)
		return nil, errDB
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetUserInfoByID(): no user record, uid: %d", uid)
		return nil, errNoUserRecord
	}

	usrinfo := &model.UserInfo{
		Name:       rec[0].Name,
		School:     rec[0].School,
		StuID:      rec[0].StuID,
		BelongTeam: rec[0].BelongTeam,
		Tshirt:     rec[0].Tshirt,
		IsUESTCStu: rec[0].IsUESTCStu,
	}

	return usrinfo, nil
}

func ModifyUserInfoByID(ctx context.Context, uid int64, usrinfo *model.UserInfoModifyReq) berrors.Berror {
	result := db.WithContext(ctx).Model(&model.UserInfoModifyReq{}).Table(tableUserInfo).Where("user_id = ?", uid).Updates(usrinfo)

	if result.Error != nil {
		hlog.Errorf("ModifyUserInfoByID(): update failed, err: %+v", result.Error)
		return errDB
	}
	if result.RowsAffected == 0 {
		hlog.Infof("ModifyUserInfoByID(): no record affected, uid: %d, info: %+v", uid, usrinfo)
	}

	return nil
}

func GetAllUsersInTeam() ([]model.User, error) {
	usrs := make([]model.User, 0)
	result := db.Model(&model.User{}).Table(tableUserInfo).Where("belong_team <> ?", 0).Find(&usrs)
	if result.Error != nil {
		return nil, result.Error
	}
	return usrs, nil
}
