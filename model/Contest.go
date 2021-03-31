package model

const (
	TableContestInfo = "t_contest"
)

type Contest struct {
	ContestID   int64  `gorm:"column:contest_id" json:"contest_id"`
	ContestName string `gorm:"column:contest_name" json:"contest_name"`
	StartTime   int64  `gorm:"column:start_time" json:"start_time"`
	EndTime     int64  `gorm:"column:end_time" json:"end_time"`
	Note        string `gorm:"column:note" json:"note"`
}
