package redis

import "github.com/HeRaNO/xcpc-team-reg/internal/berrors"

const maxIdle = 10

var (
	errRedis       = berrors.New(berrors.ErrInternal, "redis error")
	errTooFrequent = berrors.New(berrors.ErrWrongInfo, "email request too frequent")
)
