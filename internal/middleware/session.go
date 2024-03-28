package middleware

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"
)

func Authenticator() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		v := session.Get(internal.SessionName)
		if v == nil {
			c.AbortWithStatusJSON(consts.StatusOK, utils.ErrorResp(errNoSession))
			return
		}
		sessionID, ok := v.(string)
		if !ok {
			hlog.CtxErrorf(ctx, "session cannot convert to string: %+v", v)
			c.AbortWithStatusJSON(consts.StatusOK, utils.ErrorResp(errInternal))
			return
		}
		uid, err := redis.GetSession(ctx, &sessionID)
		if err != nil || uid == 0 {
			c.AbortWithStatusJSON(consts.StatusOK, utils.ErrorResp(errInternal))
			return
		}
		c.Set("uid", uid)
		c.Next(ctx)
	}
}

func CheckUnauthorized() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		v := session.Get(internal.SessionName)
		if v != nil {
			c.AbortWithStatusJSON(consts.StatusOK, utils.ErrorResp(errLoginNotExpired))
			return
		}
		c.Next(ctx)
	}
}
