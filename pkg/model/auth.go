package model

type Auth struct {
	UserID int64  `gorm:"column:user_id" json:"uid"`
	Email  string `gorm:"column:email" json:"email"`
	Pwd    string `gorm:"column:pwd" json:"pwd"`
}

type UserLoginReq struct {
	StuID    *string `json:"stuid,omitempty"`
	Email    *string `json:"email,omitempty"`
	PwdToken string  `json:"pwd"`
}
