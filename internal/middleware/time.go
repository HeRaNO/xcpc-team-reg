package middleware

import (
	"context"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func CheckBeforeEnd() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if contest.AfterRegTime(time.Now()) {
			c.AbortWithStatusJSON(consts.StatusOK, utils.ErrorResp(errEnded))
			return
		}
		c.Next(ctx)
	}
}

func CheckAfterBegin() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if contest.BeforeRegTime(time.Now()) {
			c.AbortWithStatusJSON(consts.StatusOK, utils.ErrorResp(errNotStart))
			return
		}
		c.Next(ctx)
	}
}
