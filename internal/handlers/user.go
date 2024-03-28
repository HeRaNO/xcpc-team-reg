package handlers

import (
	"context"
	"strings"

	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		hlog.Errorf("GetUserInfo(): getUID failed, err: %+v", err.Msg())
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	info, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp(info))
}

func ModifyUserInfo(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		hlog.Errorf("ModifyUserInfo(): getUID failed, err: %+v", err.Msg())
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	req := model.UserInfoModifyReq{}
	erro := c.BindAndValidate(&req)
	if erro != nil {
		hlog.Errorf("ModifyUserInfo(): BindAndValidate failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errWrongReqFmt))
		return
	}

	if req.Tshirt != nil && !contest.IsValidTshirtSize(req.Tshirt) {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidTshirtSiz))
		return
	}
	if req.School != nil && !contest.IsValidSchool(*req.School) {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidSchoolID))
		return
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			c.JSON(consts.StatusOK, utils.ErrorResp(errEmptyName))
			return
		}
		req.Name = &name
	}

	err = rdb.ModifyUserInfoByID(ctx, uid, &req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
