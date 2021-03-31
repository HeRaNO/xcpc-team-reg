package model

import (
	"context"
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
