package handlers

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.SuccessResp("xcpc-team-reg"))
}

func GetIDSchool(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.SuccessResp(contest.GetIDSchoolMap()))
}

func getUID(c *app.RequestContext) (int64, berrors.Berror) {
	v, exist := c.Get("uid")
	if !exist {
		return 0, errInvalidCookies
	}
	uid, ok := v.(int64)
	if !ok {
		return 0, errInvalidCookies
	}
	return uid, nil
}
