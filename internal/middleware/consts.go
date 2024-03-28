package middleware

import "github.com/HeRaNO/xcpc-team-reg/internal/berrors"

const (
	hashKeyLen  = 64
	blockKeyLen = 32
)

var (
	errInternal        = berrors.New(berrors.ErrInternal, "internal error")
	errNoSession       = berrors.New(berrors.ErrUnauthorized, "no session found")
	errLoginNotExpired = berrors.New(berrors.ErrUnauthorized, "login status hasn't expired")
	errNotStart        = berrors.New(berrors.ErrOutOfTime, "registration has not started yet")
	errEnded           = berrors.New(berrors.ErrOutOfTime, "registration has ended")
)
