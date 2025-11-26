package handlers

import (
	"context"
	"strings"

	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/internal/email"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Register(ctx context.Context, c *app.RequestContext) {
	req := model.UserRegisterReq{}
	erro := c.BindAndValidate(&req)
	if erro != nil {
		hlog.Errorf("Register(): BindAndValidate failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errWrongReqFmt))
		return
	}

	if !contest.IsValidTshirtSize(&req.Tshirt) {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidTshirtSiz))
		return
	}
	if !contest.IsValidSchool(req.School) {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidSchoolID))
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(consts.StatusOK, utils.ErrorResp(errEmptyName))
		return
	}

	stuID := ""
	eMail := ""
	isUestcStu := 0

	if req.StuID != nil {
		if !contest.IsValidStuID(req.StuID) {
			c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidStuID))
			return
		}
		stuID = *req.StuID
		eMail = email.MakeStuEmail(req.StuID)
		isUestcStu = 1
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
	if uid != 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(errAlreadyRegistered))
		return
	}

	err = email.ValidateEmailToken(ctx, &eMail, &req.EmailToken, &req.Action)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}

	pwdHashed, erro := utils.HashPassword(&req.PwdToken)
	if erro != nil {
		hlog.Errorf("Register(): hash password failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	regReq := rdb.UserRegInfo{
		Name:       name,
		School:     req.School,
		Email:      eMail,
		StuID:      stuID,
		Tshirt:     req.Tshirt,
		PwdToken:   *pwdHashed,
		IsUESTCStu: isUestcStu,
	}

	if err := rdb.CreateNewUser(ctx, &regReq); err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}

	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}

func ForgotPwd(ctx context.Context, c *app.RequestContext) {
	req := model.UserResetPwdReq{}
	erro := c.BindAndValidate(&req)
	if erro != nil {
		hlog.Errorf("ForgotPwd(): BindAndValidate failed, err: %+v", erro)
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

	err = email.ValidateEmailToken(ctx, &eMail, &req.EmailToken, &req.Action)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}

	pwdHashed, erro := utils.HashPassword(&req.PwdToken)
	if erro != nil {
		hlog.Errorf("ForgotPwd(): hash password failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	err = rdb.ResetUserPwd(ctx, uid, pwdHashed)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}

	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
