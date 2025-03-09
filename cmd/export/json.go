package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/bytedance/sonic"
)

type Organizations struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FormalName string `json:"formal_name"`
	Country    string `json:"country"`
}

type Teams struct {
	ID           string   `json:"id"`
	GroupIDs     []string `json:"group_ids"`
	Name         string   `json:"name"`
	DisplayName  string   `json:"display_name"`
	Members      string   `json:"members"`
	Organization string   `json:"organization_id"`
}

type Accounts struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	TeamID   string `json:"team_id"`
}

func makeMembers(info []model.UserInfo) string {
	names := make([]string, 0)
	for _, member := range info {
		names = append(names, member.Name)
	}
	return strings.Join(names, ", ")
}

func convertToJSON(info []FullTeamInfo) error {
	orgs := make([]Organizations, 0)
	teams := make([]Teams, 0)
	accounts := make([]Accounts, 0)
	id_acc_tsv := make([]byte, 0)

	org_cnt := 1
	org_map := make(map[string]string, 0)

	for _, team := range info {
		_, ok := org_map[team.TeamAffiliation]
		if !ok {
			org_map[team.TeamAffiliation] = fmt.Sprintf("org_%d", org_cnt)
			org_cnt++
		}
	}

	for org_name, org_id := range org_map {
		orgs = append(orgs, Organizations{
			ID:         org_id,
			Name:       org_name,
			FormalName: org_name,
			Country:    "CHN",
		})
	}

	for i, team := range info {
		groupID := "participants"
		if !team.IsParticipant {
			groupID = "observers"
		}
		teams = append(teams, Teams{
			ID:           team.TeamAccount,
			GroupIDs:     []string{groupID},
			Name:         team.TeamName,
			DisplayName:  team.TeamName,
			Members:      makeMembers(team.TeamMember),
			Organization: org_map[team.TeamAffiliation],
		})
		accounts = append(accounts, Accounts{
			ID:       fmt.Sprintf("account_%03d", i+1),
			Username: team.TeamAccount,
			Password: team.TeamPassword,
			Type:     "team",
			TeamID:   team.TeamAccount,
		})
		id_acc_tsv = append(id_acc_tsv, fmt.Sprintf("%s\t%s\n", team.TeamAccount, team.TeamName)...)
	}

	orgb, err := sonic.Marshal(orgs)
	if err != nil {
		return err
	}
	err = os.WriteFile("organizations.json", orgb, 0644)
	if err != nil {
		return err
	}

	teamb, err := sonic.Marshal(teams)
	if err != nil {
		return err
	}
	err = os.WriteFile("teams.json", teamb, 0644)
	if err != nil {
		return err
	}

	accountb, err := sonic.Marshal(accounts)
	if err != nil {
		return err
	}
	err = os.WriteFile("accounts.json", accountb, 0644)
	if err != nil {
		return err
	}

	return os.WriteFile("id_teamName.tsv", id_acc_tsv, 0644)
}
