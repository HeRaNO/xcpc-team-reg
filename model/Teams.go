package model

const (
	TableTeamInfo = "t_team"
)

type Team struct {
	TeamID       int64  `gorm:"column:team_id primaryKey" json:"teamid"`
	TeamName     string `gorm:"column:team_name" json:"teamname"`
	MemberCnt    int    `gorm:"column:member_cnt" json:"memcnt"`
	TeamAccount  string `gorm:"column:team_account" json:"account"`
	TeamPassword string `gorm:"column:team_password" json:"password"`
}
