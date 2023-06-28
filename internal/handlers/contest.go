package handlers

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func GetContestInfo(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.SuccessResp(contest.ContestInfo()))
}
