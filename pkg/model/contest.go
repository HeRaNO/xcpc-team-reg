package model

type ContestInfo struct {
	Name       string   `json:"name"`
	StartTime  string   `json:"start_time"`
	EndTime    string   `json:"end_time"`
	Note       string   `json:"note"`
	TshirtSize []string `json:"tshirt_size"`
}
