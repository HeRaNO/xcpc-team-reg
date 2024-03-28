package rdb

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func CreateNewUser(ctx context.Context, usr *UserRegInfo) berrors.Berror {
	trans := db.Begin()
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
		return errDB
	}
	uid := info.UserID
	authInfo := model.Auth{
		UserID: uid,
		Email:  usr.Email,
		Pwd:    usr.PwdToken,
	}
	errAuth := AddAuthInfo(ctx, &authInfo)
	if errAuth != nil {
		trans.WithContext(ctx).Rollback()
		return errDB
	}
	if err := trans.Commit().Error; nil != err {
		hlog.Errorf("CreateNewUser(): transaction failed, err: %+v", err)
		return errDB
	}
	return nil
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
