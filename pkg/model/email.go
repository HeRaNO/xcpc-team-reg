package model

type EmailVerificationReq struct {
	StuID *string `json:"stuid,omitempty"`
	Email *string `json:"email,omitempty"`
	Type  string  `json:"type"`
}
