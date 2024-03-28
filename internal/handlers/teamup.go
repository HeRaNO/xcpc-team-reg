package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func CreateTeam(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		hlog.Errorf("CreateTeam(): getUID failed, err: %+v", err.Msg())
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	usrInfo, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	if usrInfo.BelongTeam != 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInTeam))
		return
	}

	req := model.TeamInfoModifyReq{}
	erro := c.BindAndValidate(&req)
	if erro != nil {
		hlog.Errorf("CreateTeam(): BindAndValidate failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errWrongReqFmt))
		return
	}

	if req.TeamName == nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(errEmptyTeamName))
		return
	}
	name := strings.TrimSpace(*req.TeamName)
	if !validateName(&name) {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidTeamName))
		return
	}

	if req.TeamAffiliation == nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(errEmptyAffiName))
		return
	}
	affi := strings.TrimSpace(*req.TeamAffiliation)
	if !validateName(&affi) {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidAffiName))
		return
	}

	newReq := model.TeamInfoModifyReq{
		TeamName:        &name,
		TeamAffiliation: &affi,
	}
	tid, inviteToken, err := rdb.CreateNewTeam(ctx, uid, &newReq)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
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
		hlog.Errorf("JoinTeam(): getUID failed, err: %+v", err.Msg())
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	usrInfo, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	if usrInfo.BelongTeam != 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInTeam))
		return
	}

	req := model.JoinTeamReq{}
	erro := c.BindAndValidate(&req)
	if erro != nil {
		hlog.Errorf("CreateTeam(): BindAndValidate failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errWrongReqFmt))
		return
	}

	tid, erro := strconv.ParseInt(req.TeamID, 10, 64)
	if erro != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidTeamID))
		return
	}

	token, err := rdb.GetTeamInviteTokenByTeamID(ctx, tid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	if *token != req.InviteToken {
		c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidInviteToken))
		return
	}

	joined, err := rdb.UserJoinTeam(ctx, uid, tid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	if !joined {
		c.JSON(consts.StatusOK, utils.ErrorResp(errCannotJoin))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}

func QuitTeam(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		hlog.Errorf("QuitTeam(): getUID failed, err: %+v", err.Msg())
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	err = rdb.UserQuitTeam(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
