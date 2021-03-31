package modules

import (
	"net/http"
)

func ImportTeamInfo(w http.ResponseWriter, r *http.Request) {
	// read the file
	// team.team_account <- account, team.team_password <- password
}

func ExportTeamInfo(w http.ResponseWriter, r *http.Request) {
	// read database
	// data -> csv, download
}
