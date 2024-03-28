package rdb

import "github.com/HeRaNO/xcpc-team-reg/internal/berrors"

const (
	tableAuthInfo = "t_auth"
	tableTeamInfo = "t_team"
	tableUserInfo = "t_user"
)

var (
	errDB             = berrors.New(berrors.ErrInternal, "database error")
	errDuplicatedName = berrors.New(berrors.ErrWrongInfo, "team name has already been used")
	errNoTeamRecord   = berrors.New(berrors.ErrWrongInfo, "no team record")
	errNoUserRecord   = berrors.New(berrors.ErrWrongInfo, "no user record")
	errNotInATeam     = berrors.New(berrors.ErrWrongInfo, "user is not in a team")
)
