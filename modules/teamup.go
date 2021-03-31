package modules

import "net/http"

// User create a team
func CreateTeam(w http.ResponseWriter, r *http.Request) {
	// Fetch user team status
	// If user join in a team -> failed, must quit the original team
	// insert record to RDB
	// user.team_id <- team_id
}

// User join in a team
func JoinTeam(w http.ResponseWriter, r *http.Request) {
	// Fetch user team status
	// if user.team_id != 0 -> failed, has joined a team
	// Fetch team status from RDB
	// failed -> failed, not exist
	// If mem_cnt > 3 -> failed, cannot join
	// user.team_id <- team_id, team.mem_cnt++
}

// User quit a team
func QuitTeam(w http.ResponseWriter, r *http.Request) {
	// Fetch user team status
	// if user.team_id = 0 -> failed, hasn't join a team
	// Fetch team status from RDB
	// failed -> failed, not exist
	// user.team_id <- 0
	// team.mem_cnt = 1 -> delete the team record
	// else team.mem_cnt--
}
