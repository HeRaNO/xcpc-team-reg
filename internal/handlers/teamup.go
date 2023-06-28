package handlers

import (
	"context"
	"fmt"
	"html/template"
	"strconv"
	"unicode/utf8"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func CreateTeam(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	usrInfo, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	if usrInfo.BelongTeam != 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "user has already joined in a team"))
		return
	}

	req := model.TeamInfoModifyReq{}
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, err.Error()))
		return
	}

	name := ""
	if req.TeamName == nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "team name should not be empty"))
		return
	}
	name = template.HTMLEscapeString(*req.TeamName)
	nameLength := utf8.RuneCountInString(name)
	if nameLength < 1 || nameLength > contest.MaxTeamNameLength {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "team name too short or too long"))
		return
	}

	affi := ""
	if req.TeamAffiliation == nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "team affiliation should not be empty"))
		return
	}
	affi = template.HTMLEscapeString(*req.TeamAffiliation)
	affiLength := utf8.RuneCountInString(affi)
	if affiLength < 1 || affiLength > contest.MaxTeamNameLength {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "team affiliation too short or too long"))
		return
	}

	isUsed, err := redis.IsTeamNameUsed(ctx, &name)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	if isUsed {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "team name is used"))
		return
	}

	newReq := model.TeamInfoModifyReq{
		TeamName:        &name,
		TeamAffiliation: &affi,
	}
	tid, inviteToken, err := rdb.CreateNewTeam(ctx, uid, &newReq)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	resp := model.JoinTeamReq{
		TeamID:      fmt.Sprintf("%d", tid),
		InviteToken: inviteToken,
	}
	c.JSON(consts.StatusOK, utils.SuccessResp(&resp))
}

func JoinTeam(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	usrInfo, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	if usrInfo.BelongTeam != 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "user has already joined in a team"))
		return
	}

	req := model.JoinTeamReq{}
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, err.Error()))
		return
	}

	tid, err := strconv.ParseInt(req.TeamID, 10, 64)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, err.Error()))
		return
	}

	token, err := rdb.GetTeamInviteTokenByTeamID(ctx, tid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	if *token != req.InviteToken {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "invalid invite token or team id"))
		return
	}

	joined, err := rdb.UserJoinTeam(ctx, uid, tid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	if !joined {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "cannot join in the team"))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}

func QuitTeam(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	err = rdb.UserQuitTeam(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
