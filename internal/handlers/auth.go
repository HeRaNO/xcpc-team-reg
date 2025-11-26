package handlers

import (
	"context"
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
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

func validateAuthInfo(ctx context.Context, uid int64, eMail, pwdToken *string) berrors.Berror {
	info, err := rdb.GetAuthInfo(ctx, uid)
	if err != nil {
		return err
	}
	if info.Email != *eMail {
		hlog.Errorf("validateAuthInfo(): user_id in redis is different from it in rdb, uid: %d", uid)
		return errDataInconsistent
	}
	if !utils.ValidatePassword(&info.Pwd, pwdToken) {
		return errWrongPasswd
	}
	return nil
}

func Login(ctx context.Context, c *app.RequestContext) {
	req := model.UserLoginReq{}
	erro := c.BindAndValidate(&req)
	if erro != nil {
		hlog.Errorf("Login(): BindAndValidate failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errWrongReqFmt))
		return
	}

	eMail := ""
	if req.StuID != nil {
		eMail = email.MakeStuEmail(req.StuID)
	} else if req.Email != nil {
		eMail = *req.Email
	} else {
		c.JSON(consts.StatusOK, utils.ErrorResp(errNoMethod))
		return
	}

	uid, err := redis.GetUserIDByEmail(ctx, &eMail)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	if uid == 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(errNoUserRec))
		return
	}

	err = validateAuthInfo(ctx, uid, &eMail, &req.PwdToken)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}

	sid := uuid.NewString()
	err = redis.SetSession(ctx, &sid, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	session := sessions.Default(c)
	session.Set(internal.SessionName, sid)
	erro = session.Save()
	if erro != nil {
		hlog.Errorf("Login(): save session failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}

func Logout(ctx context.Context, c *app.RequestContext) {
	session := sessions.Default(c)
	sid, ok := session.Get(internal.SessionName).(string)
	if !ok {
		hlog.Errorf("Logout(): no id in session")
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidCookies))
		return
	}
	err := redis.DelSession(ctx, &sid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	session.Options(sessions.Options{
		Path:     "/",
		Domain:   internal.Domain,
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	if err := session.Save(); err != nil {
		hlog.Errorf("Logout(): save session failed, err: %+v", err)
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
