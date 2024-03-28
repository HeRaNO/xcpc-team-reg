package handlers

import "github.com/HeRaNO/xcpc-team-reg/internal/berrors"

var (
	errDataInconsistent   = berrors.New(berrors.ErrInternal, "data inconsistent")
	errInternal           = berrors.New(berrors.ErrInternal, "internal error")
	errInvalidCookies     = berrors.New(berrors.ErrInternal, "invalid cookies")
	errWrongPasswd        = berrors.New(berrors.ErrWrongInfo, "wrong password")
	errWrongReqFmt        = berrors.New(berrors.ErrWrongInfo, "wrong request format")
	errNoMethod           = berrors.New(berrors.ErrWrongInfo, "should choose one method")
	errNoUserRec          = berrors.New(berrors.ErrWrongInfo, "no such user")
	errInvalidTshirtSiz   = berrors.New(berrors.ErrWrongInfo, "invalid t-shirt size")
	errInvalidSchoolID    = berrors.New(berrors.ErrWrongInfo, "invalid school id")
	errEmptyName          = berrors.New(berrors.ErrWrongInfo, "name cannot be empty")
	errInvalidStuID       = berrors.New(berrors.ErrWrongInfo, "invalid student id")
	errAlreadyRegistered  = berrors.New(berrors.ErrWrongInfo, "email has already been registered")
	errNotInTeam          = berrors.New(berrors.ErrWrongInfo, "user hasn't joined a team")
	errInTeam             = berrors.New(berrors.ErrWrongInfo, "user has already joined in a team")
	errEmptyTeamName      = berrors.New(berrors.ErrWrongInfo, "team name should not be empty")
	errInvalidTeamName    = berrors.New(berrors.ErrWrongInfo, "team name too long or too short")
	errEmptyAffiName      = berrors.New(berrors.ErrWrongInfo, "affiliation name should not be empty")
	errInvalidAffiName    = berrors.New(berrors.ErrWrongInfo, "affiliation name too long or too short")
	errInvalidTeamID      = berrors.New(berrors.ErrWrongInfo, "invalid team id")
	errInvalidInviteToken = berrors.New(berrors.ErrWrongInfo, "invalid invite token")
	errCannotJoin         = berrors.New(berrors.ErrWrongInfo, "cannot join in the team")
)
