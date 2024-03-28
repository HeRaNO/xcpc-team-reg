package main

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/xuri/excelize/v2"
)

var heading = []string{
	"队伍编号",
	"队伍名称",
	"队伍组织",
	"是否为本校队伍",
	"队员 1 姓名",
	"队员 1 学院",
	"队员 1 学号",
	"队员 1 T 恤尺码",
	"队员 2 姓名",
	"队员 2 学院",
	"队员 2 学号",
	"队员 2 T 恤尺码",
	"队员 3 姓名",
	"队员 3 学院",
	"队员 3 学号",
	"队员 3 T 恤尺码",
}

const column = "ABCDEFGHIJKLMNOP"

func makeCellPos(i, j int) string {
	return fmt.Sprintf("%c%d", column[j], i+1)
}

func convertToXLSX(info []FullTeamInfo) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			hlog.Fatalf("xlsx close fail, err: %s", err.Error())
		}
	}()
	f.SetDocProps(&excelize.DocProperties{
		Creator: "UESTC ACM-ICPC Training Team",
	})

	teamTable := make([][]string, 0)
	teamTable = append(teamTable, heading)
	for _, team := range info {
		teamInfo := make([]string, 0)
		teamInfo = append(teamInfo, fmt.Sprintf("%d", team.TeamID))
		teamInfo = append(teamInfo, team.TeamName)
		teamInfo = append(teamInfo, team.TeamAffiliation)
		isUESTCTeam := "是"
		if !team.IsParticipant {
			isUESTCTeam = "否"
		}
		teamInfo = append(teamInfo, isUESTCTeam)
		for _, member := range team.TeamMember {
			teamInfo = append(teamInfo, member.Name)
			teamInfo = append(teamInfo, idSchoolMap[member.School])
			teamInfo = append(teamInfo, member.StuID)
			teamInfo = append(teamInfo, member.Tshirt)
		}
		teamTable = append(teamTable, teamInfo)
	}

	for i, row := range teamTable {
		for j, c := range row {
			err := f.SetCellStr("Sheet1", makeCellPos(i, j), c)
			if err != nil {
				hlog.Fatalf("set cell content failed, err: %s", err.Error())
			}
		}
	}

	err := f.SaveAs("team_export.xlsx")
	if err != nil {
		hlog.Fatalf("save file failed, err: %s", err.Error())
	}
	return nil
}
