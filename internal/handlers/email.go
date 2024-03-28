package handlers

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/email"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func SendValidationEmail(ctx context.Context, c *app.RequestContext) {
	req := model.EmailVerificationReq{}
	erro := c.BindAndValidate(&req)
	if erro != nil {
		hlog.Errorf("SendValidationEmail(): BindAndValidate failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errWrongReqFmt))
		return
	}

	e_mail := ""
	if req.StuID != nil {
		if !contest.IsValidStuID(req.StuID) {
			c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidStuID))
			return
		}
		e_mail = email.MakeStuEmail(req.StuID)
	} else if req.Email != nil {
		e_mail = *req.Email
	} else {
		c.JSON(consts.StatusOK, utils.ErrorResp(errNoMethod))
		return
	}

	err := email.SendEmailWithToken(ctx, &e_mail, &req.Type)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
