package handlers

import (
	"context"
	"html/template"
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

func GetTeamInfo(ctx context.Context, c *app.RequestContext) {
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
	if usrInfo.BelongTeam == 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "user hasn't joined a team"))
		return
	}

	info, err := rdb.GetTeamInfoByTeamID(ctx, usrInfo.BelongTeam)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, utils.SuccessResp(info))
}

func ModifyTeamInfo(ctx context.Context, c *app.RequestContext) {
	uid, err := getUID(c)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}

	req := model.TeamInfoModifyReq{}
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, err.Error()))
		return
	}

	usrInfo, err := rdb.GetUserInfoByID(ctx, uid)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	if usrInfo.BelongTeam == 0 {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "user hasn't joined a team"))
		return
	}

	if req.TeamName != nil {
		name := template.HTMLEscapeString(*req.TeamName)
		nameLength := utf8.RuneCountInString(name)
		if nameLength < 1 || nameLength > contest.MaxTeamNameLength {
			c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "team name too long or too short"))
			return
		}
		oriName, err := rdb.GetTeamNameByTeamID(ctx, usrInfo.BelongTeam)
		if err != nil {
			c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
			return
		}
		if *oriName != name {
			isUsed, err := redis.IsTeamNameUsed(ctx, &name)
			if err != nil {
				c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
				return
			}
			if isUsed {
				c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "team name is used"))
				return
			}
			req.TeamName = &name
		}
	}
	if req.TeamAffiliation != nil {
		affi := template.HTMLEscapeString(*req.TeamAffiliation)
		affiLength := utf8.RuneCountInString(affi)
		if affiLength < 1 || affiLength > contest.MaxTeamNameLength {
			c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrWrongInfo, "affiliation name too long or too short"))
			return
		}
		req.TeamAffiliation = &affi
	}

	err = rdb.ModifyTeamInfoByTeamID(ctx, usrInfo.BelongTeam, &req)
	if err != nil {
		c.JSON(consts.StatusOK, utils.ErrorResp(internal.ErrInternal, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, utils.SuccessResp("ok"))
}
