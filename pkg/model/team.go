package model

type Team struct {
	TeamID          int64  `gorm:"column:team_id;primaryKey" json:"team_id"`
	TeamName        string `gorm:"column:team_name" json:"team_name"`
	MemberCnt       int    `gorm:"column:member_cnt" json:"memcnt"`
	TeamAccount     string `gorm:"column:team_account" json:"account"`
	TeamPassword    string `gorm:"column:team_password" json:"password"`
	TeamAffiliation string `gorm:"column:team_affiliation" json:"affiliation"`
	InviteToken     string `gorm:"column:invite_token" json:"invite_token"`
}

type TeamInfo struct {
	TeamID          int64      `json:"team_id"`
	TeamName        string     `json:"team_name"`
	TeamAccount     string     `json:"account"`
	TeamPassword    string     `json:"password"`
	InviteToken     string     `json:"invite_token"`
	TeamMember      []UserInfo `json:"member"`
	MemberCnt       int        `json:"mem_cnt"`
	TeamAffiliation string     `json:"affiliation"`
}

type TeamInfoModifyReq struct {
	TeamName        *string `gorm:"column:team_name" json:"team_name,omitempty"`
	TeamAffiliation *string `gorm:"column:team_affiliation" json:"team_affiliation,omitempty"`
}

type JoinTeamReq struct {
	TeamID      string `json:"team_id"`
	InviteToken string `json:"invite_token"`
}
