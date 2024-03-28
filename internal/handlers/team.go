package handlers

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func validateName(name *string) bool {
	nameLength := utf8.RuneCountInString(*name)
	return (nameLength >= 1 && nameLength <= contest.MaxNameLength)
}

func GetTeamInfo(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		hlog.Errorf("GetTeamInfo(): getUID failed, err: %+v", err.Msg())
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	usrInfo, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	if usrInfo.BelongTeam == 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(errNotInTeam))
		return
	}

	info, err := rdb.GetTeamInfoByTeamID(ctx, usrInfo.BelongTeam)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}

	c.JSON(consts.StatusOK, utils.SuccessResp(info))
}

func ModifyTeamInfo(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		hlog.Errorf("ModifyTeamInfo(): getUID failed, err: %+v", err.Msg())
		c.JSON(consts.StatusOK, utils.ErrorResp(errInternal))
		return
	}

	req := model.TeamInfoModifyReq{}
	erro := c.BindAndValidate(&req)
	if erro != nil {
		hlog.Errorf("ModifyTeamInfo(): BindAndValidate failed, err: %+v", erro)
		c.JSON(consts.StatusOK, utils.ErrorResp(errWrongReqFmt))
		return
	}

	usrInfo, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	if usrInfo.BelongTeam == 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(errNotInTeam))
		return
	}

	if req.TeamName != nil {
		teamNameTrimed := strings.TrimSpace(*req.TeamName)
		if !validateName(&teamNameTrimed) {
			c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidTeamName))
			return
		}
		req.TeamName = &teamNameTrimed
	}
	if req.TeamAffiliation != nil {
		teamAffiTrimed := strings.TrimSpace(*req.TeamAffiliation)
		if !validateName(&teamAffiTrimed) {
			c.JSON(consts.StatusOK, utils.ErrorResp(errInvalidAffiName))
			return
		}
		req.TeamAffiliation = &teamAffiTrimed
	}

	err = rdb.ModifyTeamInfoByTeamID(ctx, usrInfo.BelongTeam, &req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(err))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
