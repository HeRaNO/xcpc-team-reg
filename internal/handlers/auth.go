package handlers

import (
	"context"
	"errors"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/internal/email"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"github.com/hertz-contrib/sessions"
)

func validateAuthInfo(ctx context.Context, uid int64, e_mail *string, pwdToken *string) (bool, error) {
	info, err := rdb.GetAuthInfo(ctx, uid)
	if err != nil {
		return false, err
	}
	if info.Email != *e_mail {
		hlog.Errorf("validateAuthInfo(): user_id in Redis is different from it in rdb")
		return false, errors.New("data inconsistent")
	}
	if !utils.ValidatePassword(&info.Pwd, pwdToken) {
		return true, errors.New("wrong password")
	}
	return false, nil
}

func Login(ctx context.Context, c *app.RequestContext) {
	req := model.UserLoginReq{}
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.Errorf("Login(): BindAndValidate failed, err: %+v", err)
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, err.Error()))
		return
	}

	e_mail := ""
	if req.StuID != nil {
		e_mail = email.MakeStuEmail(req.StuID)
	} else if req.Email != nil {
		e_mail = *req.Email
	} else {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "should choose one method to verify email"))
		return
	}

	uid, err := redis.GetUserIDByEmail(ctx, &e_mail)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	if uid == 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "no such user"))
		return
	}

	flag, err := validateAuthInfo(ctx, uid, &e_mail, &req.PwdToken)
	if err != nil {
		errCode := internal.ErrInternal
		if flag {
			errCode = internal.ErrWrongInfo
		}
		c.JSON(consts.StatusOK, utils.ErrorResp(errCode, err.Error()))
		return
	}

	sid := uuid.NewString()
	err = redis.SetSession(ctx, &sid, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	session := sessions.Default(c)
	session.Set(internal.SessionName, sid)
	err = session.Save()
	if err != nil {
		hlog.Errorf("Login(): save session failed, err: %+v", err)
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}

func Logout(ctx context.Context, c *app.RequestContext) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		hlog.Errorf("Logout(): save session failed, err: %+v", err)
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
