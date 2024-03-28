package main

import (
	"flag"
	"fmt"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func main() {
	initConfigFile := flag.String("c", "./configs/conf.yaml", "the path of configure file")
	flag.Parse()
	internal.InitConfig(initConfigFile)

	teams, err := rdb.GetAllTeams()
	if err != nil {
		hlog.Fatalf("cannot get all teams, err: %s", err.Error())
	}

	for i, team := range teams {
		teamName := fmt.Sprintf("team%03d", i+1)
		teamPwd, err := utils.GenToken(contest.UserTokenLength)
		if err != nil {
			hlog.Fatalf("cannot gen password, err: %s", err.Msg())
		}
		erro := rdb.SetTeamAccPwdByID(team.TeamID, &teamName, &teamPwd)
		if erro != nil {
			hlog.Fatalf("update failed, err: %s", erro.Error())
		}
	}

	hlog.Info("gen account finished successfully")
}
