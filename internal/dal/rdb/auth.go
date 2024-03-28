package rdb

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func AddAuthInfo(ctx context.Context, info *model.Auth) berrors.Berror {
	err := redis.SetEmailUserID(ctx, info.UserID, &info.Email)
	if err != nil {
		return err
	}
	erro := db.WithContext(ctx).Model(&model.Auth{}).Table(tableAuthInfo).Create(info).Error
	if erro != nil {
		hlog.Errorf("AddAuthInfo(): create failed, err: %+v", erro)
		redis.DelEmailUserID(ctx, &info.Email)
		return errDB
	}
	return nil
}

func GetAuthInfo(ctx context.Context, uid int64) (*model.Auth, berrors.Berror) {
	rec := make([]model.Auth, 0)
	result := db.WithContext(ctx).Model(&model.Auth{}).Table(tableAuthInfo).Where("user_id = ?", uid).Find(&rec)

	if result.Error != nil {
		hlog.Errorf("GetAuthInfo(): query failed, err: %+v", result.Error)
		return nil, errDB
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetAuthInfo(): no user record, uid: %d", uid)
		return nil, errDB
	}
	if result.RowsAffected > 1 {
		hlog.Errorf("GetAuthInfo(): duplicate uid in t_auth, uid: %d", uid)
		return nil, errDB
	}

	return &rec[0], nil
}

func ResetUserPwd(ctx context.Context, uid int64, pwdToken *string) berrors.Berror {
	err := db.WithContext(ctx).Model(&model.Auth{}).Table(tableAuthInfo).Where("user_id = ?", uid).Update("pwd", *pwdToken).Error
	if err != nil {
		hlog.Errorf("ResetUserPwd(): update failed, err: %+v", err)
		return errDB
	}
	return nil
}
