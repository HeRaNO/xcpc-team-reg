package modules

import (
	"net/http"
)

func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func ImportTeamInfo(w http.ResponseWriter, r *http.Request) {
	// read the file
	// team.team_account <- account, team.team_password <- password
}

func ExportTeamInfo(w http.ResponseWriter, r *http.Request) {
	// read database
	// data -> csv, download
}
