package handlers

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/email"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func SendValidationEmail(ctx context.Context, c *app.RequestContext) {
	req := model.EmailVerificationReq{}
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.Errorf("SendValidationEmail(): BindAndValidate failed, err: %+v", err)
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

	flag, err := email.SendEmailWithToken(ctx, &e_mail, &req.Type)
	if err != nil {
		errCode := internal.ErrInternal
		if flag {
			errCode = internal.ErrWrongInfo
		}
		c.JSON(consts.StatusOK, utils.ErrorResp(errCode, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
