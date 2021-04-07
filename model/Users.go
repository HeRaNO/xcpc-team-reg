package model

import (
	"context"
	"errors"
	"log"

	"github.com/HeRaNO/xcpc-team-reg/config"
)

const (
	TableUserInfo = "t_user"
)

type UserRegister struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	School     int    `json:"school"`
	StuID      string `json:"stuid"`
	EmailToken string `json:"email_token"`
	PwdToken   string `json:"pwd_token"`
	Action     string `json:"action"`
}

type UserInfo struct {
	Name       string `json:"name"`
	School     string `json:"school"`
	StuID      string `json:"stuid"`
	BelongTeam int64  `json:"teamid"`
}

type UserInfoModify struct {
	Name   string `gorm:"column:user_name" json:"name"`
	School int    `gorm:"column:school" json:"school"`
	StuID  string `gorm:"column:stu_id" json:"stuid"`
}

type User struct {
	UserID     int64  `gorm:"column:user_id;primaryKey" json:"userid"`
	Name       string `gorm:"column:user_name" json:"name"`
	Email      string `gorm:"column:email" json:"email"`
	School     int    `gorm:"column:school" json:"school"`
	StuID      string `gorm:"column:stu_id" json:"stuid"`
	BelongTeam int64  `gorm:"column:belong_team" json:"teamid"`
	IsAdmin    int    `gorm:"column:is_admin" json:"is_admin"`
}

func CreateNewUser(ctx context.Context, usr UserRegister, isAdmin int) error {
	// INSERT INTO user(...) VALUES [values in usr]
	trans := config.RDB.Begin()
	info := User{
		Name:       usr.Name,
		Email:      usr.Email,
		School:     usr.School,
		StuID:      usr.StuID,
		BelongTeam: 0,
		IsAdmin:    isAdmin,
	}
	err := trans.WithContext(ctx).Model(&User{}).Table(TableUserInfo).Create(&info).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}
	uid := info.UserID
	authInfo := Auth{
		UserID: uid,
		Email:  usr.Email,
		Pwd:    usr.PwdToken,
	}
	err = AddAuthInfo(ctx, &authInfo)
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}
	if err := trans.Commit().Error; nil != err {
		log.Println("[ERROR] CreateNewUser(): transaction failed")
		return err
	}
	return nil
}

func GetUserInfoByID(ctx context.Context, uid int64) (*UserInfo, error) {
	rdb := config.RDB

	rec := make([]User, 0)
	result := rdb.Model(&User{}).Table(TableUserInfo).Where("user_id = ?", uid).Find(&rec)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no user record")
	}

	if result.RowsAffected > 1 {
		return nil, errors.New("duplicate user_id but why???")
	}

	usrSchool := "undefined"

	if school, ok := config.SchoolMap[rec[0].School]; ok {
		usrSchool = school
	}

	usrinfo := &UserInfo{
		Name:       rec[0].Name,
		School:     usrSchool,
		StuID:      rec[0].StuID,
		BelongTeam: rec[0].BelongTeam,
	}

	return usrinfo, nil
}

func GetAdminByUserID(ctx context.Context, uid int64) (bool, error) {
	rdb := config.RDB

	rec := map[string]interface{}{}
	result := rdb.Model(&User{}).Table(TableUserInfo).Select("is_admin").Where("user_id = ?", uid).Find(&rec)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, errors.New("no user record")
	}

	if result.RowsAffected > 1 {
		return false, errors.New("duplicate user_id but why???")
	}

	isAdmin := rec["is_admin"].(int)

	if isAdmin == 0 {
		return false, nil
	}

	if isAdmin != 1 {
		return false, errors.New("why the is_admin is not 0 or 1")
	}

	return true, nil
}

func GetTeamIDByUserID(ctx context.Context, uid int64) (int64, error) {
	rdb := config.RDB

	rec := map[string]interface{}{}
	result := rdb.Model(&User{}).Table(TableUserInfo).Select("belong_team").Where("user_id = ?", uid).Find(&rec)

	if result.Error != nil {
		return -1, result.Error
	}

	if result.RowsAffected == 0 {
		return -1, errors.New("no user record")
	}

	if result.RowsAffected > 1 {
		return -1, errors.New("duplicate user_id but why???")
	}

	teamID := rec["belong_team"].(int64)

	return teamID, nil
}

func ModifyUserInfoByID(ctx context.Context, uid int64, usrinfo *UserInfoModify) error {
	trans := config.RDB.Begin()

	err := trans.WithContext(ctx).Model(&UserInfoModify{}).Table(TableUserInfo).Where("user_id = ?", uid).Updates(usrinfo).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}

	if err := trans.Commit().Error; err != nil {
		log.Println("[ERROR] ModifyUserInfoByID(): transaction failed")
		return err
	}

	return nil
}

func GetUserInfosByTeamID(ctx context.Context, tid int64) ([]UserInfo, error) {
	rdb := config.RDB

	rec := make([]User, config.MaxTeamMember)
	result := rdb.Model(&User{}).Table(TableUserInfo).Where("belong_team = ?", tid).Find(&rec)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no user in this team but why???")
	}

	if result.RowsAffected > int64(config.MaxTeamMember) {
		return nil, errors.New("too many members in this team but why???")
	}

	usrInfo := make([]UserInfo, config.MaxTeamMember)

	for _, usr := range rec {
		usrSchool := "undefined"

		if school, ok := config.SchoolMap[usr.School]; ok {
			usrSchool = school
		}
		info := UserInfo{
			Name:       usr.Name,
			School:     usrSchool,
			StuID:      usr.StuID,
			BelongTeam: usr.BelongTeam,
		}

		usrInfo = append(usrInfo, info)
	}

	return usrInfo, nil
}
