package main

import (
	"context"
	"flag"
	"strconv"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type FullTeamInfo struct {
	TeamName        string
	TeamAccount     string
	TeamPassword    string
	TeamMember      []model.UserInfo
	TeamAffiliation string
	IsParticipant   bool
	IsFreshman      bool
}

var idSchoolMap map[int]string

func isFreshman(id string) bool {
	if len(id) != 13 {
		return false
	}
	currentYearStr := strconv.Itoa(time.Now().Year() - 1)
	return id[:4] == currentYearStr
}

func main() {
	initConfigFile := flag.String("c", "./configs/conf.yaml", "the path of configure file")
	exportForm := flag.String("f", "xlsx", "export form of team info, can be: xlsx and json")
	flag.Parse()
	internal.InitConfig(initConfigFile)

	ctx := context.Background()
	idSchoolMap = contest.GetIDSchoolMap()

	teamInfos, err := rdb.GetAllTeamInfos(ctx)
	if err != nil {
		hlog.Fatalf("cannot get all team infos, err: %s", err.Error())
	}

	allTeamInfo := make([]FullTeamInfo, 0, len(teamInfos))

	for _, fullTeamInfo := range teamInfos {
		isUESTCTeam := true
		isFreshmanTeam := false
		for _, member := range fullTeamInfo.TeamMember {
			if member.IsUESTCStu != 1 {
				isUESTCTeam = false
				break
			}
		}
		if isUESTCTeam {
			fullTeamInfo.TeamAffiliation = "UESTC"
			isFreshmanTeam = true
			for _, member := range fullTeamInfo.TeamMember {
				if !isFreshman(member.StuID) {
					isFreshmanTeam = false
					break
				}
			}
		}
		teamInfo := FullTeamInfo{
			TeamName:        fullTeamInfo.TeamName,
			TeamAccount:     fullTeamInfo.TeamAccount,
			TeamPassword:    fullTeamInfo.TeamPassword,
			TeamMember:      fullTeamInfo.TeamMember,
			TeamAffiliation: fullTeamInfo.TeamAffiliation,
			IsParticipant:   isUESTCTeam,
			IsFreshman:      isFreshmanTeam,
		}
		allTeamInfo = append(allTeamInfo, teamInfo)
	}

	switch *exportForm {
	case "xlsx":
		err := convertToXLSX(allTeamInfo)
		if err != nil {
			hlog.Fatalf("cannot gen xlsx file, err: %s", err.Error())
		}
	case "json":
		err := convertToJSON(allTeamInfo)
		if err != nil {
			hlog.Fatalf("cannot gen json file, err: %s", err.Error())
		}
	default:
		hlog.Fatalf("unexpected form: %s", *exportForm)
	}

	hlog.Info("export file(s) successfully")
}
