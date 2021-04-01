package modules

import "net/http"

func GetTeamInfo(w http.ResponseWriter, r *http.Request) {
	// Fetch from RDB
	// return info

}

func ModifyTeamInfo(w http.ResponseWriter, r *http.Request) {
	// validate info
	// write to RDB
}
