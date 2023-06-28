package handlers

import (
	"context"
	"html/template"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	info, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp(info))
}

func ModifyUserInfo(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	req := model.UserInfoModifyReq{}
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, err.Error()))
		return
	}

	if req.Tshirt != nil && !contest.IsValidTshirtSize(req.Tshirt) {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "invalid T-shirt size"))
		return
	}
	if req.School != nil && !contest.IsValidSchool(*req.School) {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "invalid school id"))
		return
	}
	if req.Name != nil {
		name := template.HTMLEscapeString(*req.Name)
		req.Name = &name
	}

	err = rdb.ModifyUserInfoByID(ctx, uid, &req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
