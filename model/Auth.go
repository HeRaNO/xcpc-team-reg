package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/go-redis/redis/v8"
)

const (
	TableAuthInfo = "t_auth"
)

type Auth struct {
	UserID int64  `gorm:"column:user_id" json:"uid"`
	Email  string `gorm:"column:email" json:"email"`
	Pwd    string `gorm:"column:pwd" json:"pwd"`
}

type EmailVerification struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

type UserLogin struct {
	Email    string `json:"email"`
	PwdToken string `json:"pwd"`
}

func MakeEmailTokenKey(email *string) string {
	ret := fmt.Sprintf("EMAILTOKEN:%s", *email)
	return ret
}

func MakeEmailUserIDKey(email *string) string {
	ret := fmt.Sprintf("EMAIL:%s", *email)
	return ret
}

func MakeEmailRequestKey(email *string) string {
	ret := fmt.Sprintf("EMAILREQ:%s", *email)
	return ret
}

func MakeEmailActionKey(email *string) string {
	ret := fmt.Sprintf("EMAILACTION:%s", *email)
	return ret
}

func GetEmailToken(ctx context.Context, email *string) (string, error) {
	key := MakeEmailTokenKey(email)
	ret, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[INFO] GetEmailToken(): key is nil, email: %s\n", *email)
		return "", nil
	} else if err != nil {
		log.Println("[ERROR] GetEmailToken(): redis query error")
		return "", err
	}
	return ret, nil
}

func SetEmailToken(ctx context.Context, email *string, token *string) error {
	key := MakeEmailTokenKey(email)
	err := config.RedisClient.SetEX(ctx, key, *token, config.EMAILTOKEN_EXPIRETIME).Err()
	if err != nil {
		log.Println("[ERROR] SetEmailToken(): redis set error")
		return err
	}
	return nil
}

func DelEmailToken(ctx context.Context, email *string) error {
	key := MakeEmailTokenKey(email)
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		log.Println("[ERROR] DelEmailToken(): redis del error")
		return err
	}
	return nil
}

func GetUserIDByEmail(ctx context.Context, email *string) (int64, error) {
	var uid int64
	key := MakeEmailUserIDKey(email)
	ret, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[INFO] GetUserIDByEmail(): key is nil, email: %s\n", *email)
		return -1, nil
	} else if err != nil {
		log.Println("[ERROR] GetUserIDByEmail(): redis query error")
		return -1, err
	}
	uid, err = strconv.ParseInt(ret, 10, 64)
	if err != nil {
		log.Printf("[ERROR] GetUserIDByEmail(): parse uid failed, uid from Redis: %s", ret)
		return -1, err
	}
	return uid, nil
}

func SetEmailUserID(ctx context.Context, uid int64, email *string) error {
	key := MakeEmailUserIDKey(email)
	err := config.RedisClient.Set(ctx, key, uid, 0).Err()
	if err != nil {
		log.Println("[ERROR] SetEmailUserID(): redis set error")
		return err
	}
	return nil
}

func DelEmailUserID(ctx context.Context, email *string) error {
	key := MakeEmailUserIDKey(email)
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		log.Println("[ERROR] DelEmailUserID(): redis del error")
		return err
	}
	return nil
}

func GetEmailRequest(ctx context.Context, email *string) error {
	key := MakeEmailRequestKey(email)
	err := config.RedisClient.Get(ctx, key).Err()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		log.Println("[ERROR] GetEmailRequest(): redis query error")
		return err
	}
	return errors.New("email request too frequent")
}

func SetEmailRequest(ctx context.Context, email *string) error {
	key := MakeEmailRequestKey(email)
	err := config.RedisClient.SetEX(ctx, key, "1", config.EMAILSEND_GAPTIME).Err()
	if err != nil {
		log.Println("[ERROR] SetEmailRequest(): redis set error")
		return err
	}
	return nil
}

func GetEmailAction(ctx context.Context, email *string) (string, error) {
	key := MakeEmailActionKey(email)
	ret, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		log.Println("[ERROR] GetEmailRequest(): redis query error")
		return "", err
	}
	return ret, nil
}

func SetEmailAction(ctx context.Context, email *string, action *string) error {
	key := MakeEmailActionKey(email)
	err := config.RedisClient.SetEX(ctx, key, *action, config.EMAILTOKEN_EXPIRETIME).Err()
	if err != nil {
		log.Println("[ERROR] SetEmailAction(): redis set error")
		return err
	}
	return nil
}

func DelEmailAction(ctx context.Context, email *string) error {
	key := MakeEmailActionKey(email)
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		log.Println("[ERROR] DelEmailAction(): redis del error")
		return err
	}
	return nil
}

func ValidateEmailToken(ctx context.Context, email *string, token *string, action *string) error {
	actionFromRedis, err := GetEmailAction(ctx, email)
	if err != nil {
		return err
	}
	if actionFromRedis == "" || actionFromRedis != *action {
		return errors.New("action is invalid")
	}
	tokenFromRedis, err := GetEmailToken(ctx, email)
	if err != nil {
		return err
	}
	if tokenFromRedis == "" || tokenFromRedis != *token {
		return errors.New("token is invalid")
	}
	err = DelEmailToken(ctx, email)
	if err != nil {
		return err
	}
	err = DelEmailAction(ctx, email)
	if err != nil {
		return err
	}
	return nil
}

func AddAuthInfo(ctx context.Context, info *Auth) error {
	trans := config.RDB.Begin()
	err := SetEmailUserID(ctx, info.UserID, &info.Email)
	if err != nil {
		return err
	}

	err = trans.WithContext(ctx).Model(&Auth{}).Table(TableAuthInfo).Create(info).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		DelEmailUserID(ctx, &info.Email)
		return err
	}
	if err := trans.Commit().Error; nil != err {
		log.Println("[ERROR] AddAuthInfo(): transaction failed")
		return err
	}
	return nil
}

func ValidateAuthInfo(ctx context.Context, uid int64, email *string, token *string) error {
	rdb := config.RDB

	rec := make([]Auth, 0)
	result := rdb.WithContext(ctx).Model(&Auth{}).Table(TableAuthInfo).Where("user_id = ?", uid).Find(&rec)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no user record")
	}

	if result.RowsAffected > 1 {
		return errors.New("duplicate user_id but why???")
	}

	usrAuth := rec[0]

	if *email != usrAuth.Email {
		log.Println("[ERROR] ValidateAuthInfo(): user_id in Redis is different from it in database")
		return errors.New("data inconsistent")
	}

	if *token != usrAuth.Pwd {
		return errors.New("wrong password")
	}

	return nil
}

func ResetUserPwd(ctx context.Context, uid int64, pwdToken *string) error {
	trans := config.RDB.Begin()

	err := trans.WithContext(ctx).Model(&Auth{}).Table(TableAuthInfo).Where("user_id = ?", uid).Update("pwd", *pwdToken).Error
	if err != nil {
		trans.WithContext(ctx).Rollback()
		return err
	}

	if err := trans.Commit().Error; nil != err {
		log.Println("[ERROR] ResetUserPwd(): transaction failed")
		return err
	}
	return nil
}
