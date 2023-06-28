package model

type User struct {
	UserID     int64  `gorm:"column:user_id;primaryKey" json:"userid"`
	Name       string `gorm:"column:user_name" json:"name"`
	Email      string `gorm:"column:email" json:"email"`
	School     int    `gorm:"column:school" json:"school"`
	StuID      string `gorm:"column:stu_id" json:"stuid"`
	BelongTeam int64  `gorm:"column:belong_team" json:"teamid"`
	IsUESTCStu int    `gorm:"column:is_uestc_stu" json:"is_uestc_stu"`
	Tshirt     string `gorm:"column:tshirt" json:"tshirt"`
}

type UserInfo struct {
	Name       string `json:"name"`
	School     int    `json:"school"`
	StuID      string `json:"stuid"`
	BelongTeam int64  `json:"teamid"`
	Tshirt     string `json:"tshirt"`
	IsUESTCStu int    `json:"is_uestc_stu"`
}

type UserRegisterReq struct {
	Name       string  `json:"name"`
	School     int     `json:"school"`
	Email      *string `json:"email,omitempty"`
	StuID      *string `json:"stuid,omitempty"`
	Tshirt     string  `json:"tshirt"`
	EmailToken string  `json:"email_token"`
	PwdToken   string  `json:"pwd_token"`
	Action     string  `json:"action"`
}

type UserResetPwdReq struct {
	StuID      *string `json:"stuid,omitempty"`
	Email      *string `json:"email,omitempty"`
	EmailToken string  `json:"email_token"`
	PwdToken   string  `json:"pwd_token"`
	Action     string  `json:"action"`
}

type UserInfoModifyReq struct {
	Name   *string `gorm:"column:user_name" json:"name,omitempty"`
	School *int    `gorm:"column:school" json:"school,omitempty"`
	Tshirt *string `gorm:"column:tshirt" json:"tshirt,omitempty"`
}
