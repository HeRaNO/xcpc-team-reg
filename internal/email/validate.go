package email

import (
	"context"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
)

func ValidateEmailToken(ctx context.Context, email *string, token *string, action *string) berrors.Berror {
	actionFromRedis, err := redis.GetEmailAction(ctx, email)
	if err != nil {
		return err
	}
	if actionFromRedis == "" || actionFromRedis != *action {
		return errInvalidType
	}
	tokenFromRedis, err := redis.GetEmailToken(ctx, email)
	if err != nil {
		return err
	}
	if tokenFromRedis == "" || tokenFromRedis != *token {
		return errInvalidToken
	}
	redis.DelEmailToken(ctx, email)
	redis.DelEmailAction(ctx, email)
	return nil
}
