package handlers

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal"
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
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, err.Error()))
		return
	}

	if !contest.IsValidTshirtSize(&req.Tshirt) {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "invalid t-shirt size"))
		return
	}
	if !contest.IsValidSchool(req.School) {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "invalid school id"))
		return
	}
	name := utils.TrimName(&req.Name)
	if *name == "" {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "name cannot be empty"))
		return
	}

	stuID := ""
	e_mail := ""
	is_uestc_stu := 0

	if req.StuID != nil {
		if !contest.IsValidStuID(req.StuID) {
			c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "invalid student id"))
			return
		}
		stuID = *req.StuID
		e_mail = email.MakeStuEmail(req.StuID)
		is_uestc_stu = 1
	} else if req.Email != nil {
		e_mail = *req.Email
	} else {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "should choose one register method"))
		return
	}

	uid, err := redis.GetUserIDByEmail(ctx, &e_mail)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	if uid != 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "email has already registered"))
		return
	}

	flag, err := email.ValidateEmailToken(ctx, &e_mail, &req.EmailToken, &req.Action)
	if err != nil {
		errCode := internal.ErrInternal
		if flag {
			errCode = internal.ErrWrongInfo
		}
		c.JSON(consts.StatusOK, utils.ErrorResp(errCode, err.Error()))
		return
	}

	pwdHashed, err := utils.HashPassword(&req.PwdToken)
	if err != nil {
		hlog.Errorf("Register(): hash password failed, err: %+v", err)
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	regReq := rdb.UserRegInfo{
		Name:       *name,
		School:     req.School,
		Email:      e_mail,
		StuID:      stuID,
		Tshirt:     req.Tshirt,
		PwdToken:   *pwdHashed,
		IsUESTCStu: is_uestc_stu,
	}

	if err := rdb.CreateNewUser(ctx, &regReq); err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}

func ForgotPwd(ctx context.Context, c *app.RequestContext) {
	req := model.UserResetPwdReq{}
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, err.Error()))
		return
	}

	e_mail := ""
	if req.StuID != nil {
		e_mail = email.MakeStuEmail(req.StuID)
	} else if req.Email != nil {
		e_mail = *req.Email
	} else {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "should choose one reset method"))
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

	flag, err := email.ValidateEmailToken(ctx, &e_mail, &req.EmailToken, &req.Action)
	if err != nil {
		errCode := internal.ErrInternal
		if flag {
			errCode = internal.ErrWrongInfo
		}
		c.JSON(consts.StatusOK, utils.ErrorResp(errCode, err.Error()))
		return
	}

	pwdHashed, err := utils.HashPassword(&req.PwdToken)
	if err != nil {
		hlog.Errorf("ForgotPwd(): hash password failed, err: %+v", err)
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	err = rdb.ResetUserPwd(ctx, uid, pwdHashed)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
