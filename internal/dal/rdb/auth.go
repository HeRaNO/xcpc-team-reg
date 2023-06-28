package rdb

import (
	"context"
	"errors"

	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func AddAuthInfo(ctx context.Context, info *model.Auth) error {
	err := redis.SetEmailUserID(ctx, info.UserID, &info.Email)
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).Model(&model.Auth{}).Table(tableAuthInfo).Create(info).Error
	if err != nil {
		hlog.Errorf("AddAuthInfo(): create failed, err: %+v", err)
		redis.DelEmailUserID(ctx, &info.Email)
		return err
	}
	return nil
}

func GetAuthInfo(ctx context.Context, uid int64) (*model.Auth, error) {
	rec := make([]model.Auth, 0)
	result := db.WithContext(ctx).Model(&model.Auth{}).Table(tableAuthInfo).Where("user_id = ?", uid).Find(&rec)

	if result.Error != nil {
		hlog.Errorf("GetAuthInfo(): query failed, err: %+v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		hlog.Infof("GetAuthInfo(): no user record, uid: %d", uid)
		return nil, errors.New("no user record")
	}
	if result.RowsAffected > 1 {
		hlog.Errorf("GetAuthInfo(): duplicate uid in t_auth, uid: %d", uid)
		return nil, errors.New("duplicate user_id but why???")
	}

	return &rec[0], nil
}

func ResetUserPwd(ctx context.Context, uid int64, pwdToken *string) error {
	err := db.WithContext(ctx).Model(&model.Auth{}).Table(tableAuthInfo).Where("user_id = ?", uid).Update("pwd", *pwdToken).Error
	if err != nil {
		hlog.Errorf("ResetUserPwd(): update failed, err: %+v", err)
		return err
	}
	return nil
}
