package main

import (
	"context"
	"flag"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type FullTeamInfo struct {
	TeamID          int64
	TeamName        string
	TeamAccount     string
	TeamPassword    string
	TeamMember      []model.UserInfo
	TeamAffiliation string
	IsParticipant   bool
}

var idSchoolMap map[int]string

func main() {
	initConfigFile := flag.String("c", "./configs/conf.yaml", "the path of configure file")
	exportForm := flag.String("f", "xlsx", "export form of team info, can be: xlsx, json and tsv")
	useExternalID := flag.Bool("external", false, "use externalID in json, import to DOMJudge")
	flag.Parse()
	internal.InitConfig(initConfigFile)

	ctx := context.Background()
	idSchoolMap = contest.GetIDSchoolMap()

	teams, err := rdb.GetAllTeams()
	if err != nil {
		hlog.Fatalf("cannot get all teams, err: %s", err.Error())
	}

	allTeamInfo := make([]FullTeamInfo, 0)

	for _, team := range teams {
		fullTeamInfo, err := rdb.GetTeamInfoByTeamID(ctx, team.TeamID)
		if err != nil {
			hlog.Fatalf("cannot get team info, tid: %d, err: %s", team.TeamID, err.Error())
		}
		isUESTCTeam := true
		for _, member := range fullTeamInfo.TeamMember {
			if member.IsUESTCStu != 1 {
				isUESTCTeam = false
				break
			}
		}
		if isUESTCTeam {
			fullTeamInfo.TeamAffiliation = "UESTC"
		}
		teamInfo := FullTeamInfo{
			TeamID:          fullTeamInfo.TeamID,
			TeamName:        fullTeamInfo.TeamName,
			TeamAccount:     fullTeamInfo.TeamAccount,
			TeamPassword:    fullTeamInfo.TeamPassword,
			TeamMember:      fullTeamInfo.TeamMember,
			TeamAffiliation: fullTeamInfo.TeamAffiliation,
			IsParticipant:   isUESTCTeam,
		}
		allTeamInfo = append(allTeamInfo, teamInfo)
	}

	if *exportForm == "xlsx" {
		err := convertToXLSX(allTeamInfo)
		if err != nil {
			hlog.Fatalf("cannot gen xlsx file, err: %s", err.Error())
		}
	} else if *exportForm == "json" {
		err := convertToJSON(allTeamInfo, *useExternalID)
		if err != nil {
			hlog.Fatalf("cannot gen json file, err: %s", err.Error())
		}
	} else if *exportForm == "tsv" {
		err := convertToTSV(allTeamInfo)
		if err != nil {
			hlog.Fatalf("cannot gen tsv file, err: %s", err.Error())
		}
	} else {
		hlog.Fatalf("unexpected form: %s", *exportForm)
	}

	hlog.Info("export file(s) successfully")
}
