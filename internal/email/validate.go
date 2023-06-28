package email

import (
	"context"
	"errors"

	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
)

func ValidateEmailToken(ctx context.Context, email *string, token *string, action *string) (bool, error) {
	actionFromRedis, err := redis.GetEmailAction(ctx, email)
	if err != nil {
		return false, err
	}
	if actionFromRedis == "" || actionFromRedis != *action {
		return true, errors.New("action is invalid")
	}
	tokenFromRedis, err := redis.GetEmailToken(ctx, email)
	if err != nil {
		return false, err
	}
	if tokenFromRedis == "" || tokenFromRedis != *token {
		return true, errors.New("token is invalid")
	}
	redis.DelEmailToken(ctx, email)
	redis.DelEmailAction(ctx, email)
	return true, nil
}
