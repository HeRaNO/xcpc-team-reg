package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
)

func writeTSV(fileName string, content [][]string) error {
	tsvContent := make([]byte, 0)
	buf := bytes.NewBuffer(tsvContent)
	tsvWriter := csv.NewWriter(buf)
	tsvWriter.Comma = '\t'

	err := tsvWriter.WriteAll(content)
	if err != nil {
		return err
	}

	tsvWriter.Flush()
	return os.WriteFile(fileName, buf.Bytes(), 0644)
}

// Deprecated: The legacy TSV importing method is rather unsafe
// for importing, and too much inconvenience. We will deprecate
// it next year.
func convertToTSV(info []FullTeamInfo) error {
	teams := make([][]string, 0)
	accounts := make([][]string, 0)

	teams = append(teams, []string{"File_Version", "2"})
	accounts = append(accounts, []string{"accounts", "1"})

	for i, team := range info {
		teamInfo := make([]string, 0)
		teamInfo = append(teamInfo, fmt.Sprintf("%d", i+1))
		teamInfo = append(teamInfo, "")
		groupID := "3"
		if !team.IsParticipant {
			groupID = "4"
		}
		teamInfo = append(teamInfo, groupID)
		teamInfo = append(teamInfo, team.TeamName)
		teamInfo = append(teamInfo, team.TeamAffiliation)
		teamInfo = append(teamInfo, team.TeamAffiliation)
		teamInfo = append(teamInfo, "CHN")
		teams = append(teams, teamInfo)

		accoutInfo := make([]string, 0)
		accoutInfo = append(accoutInfo, "team")
		accoutInfo = append(accoutInfo, team.TeamAccount)
		accoutInfo = append(accoutInfo, team.TeamAccount)
		accoutInfo = append(accoutInfo, team.TeamPassword)
		accounts = append(accounts, accoutInfo)
	}

	if err := writeTSV("teams.tsv", teams); err != nil {
		return err
	}
	return writeTSV("accounts.tsv", accounts)
}
