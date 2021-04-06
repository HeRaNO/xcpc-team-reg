package modules

import (
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/util"
)

func GetContestInfo(w http.ResponseWriter, r *http.Request) {
	// return info
	info := map[string]string{
		"name":       config.ContestName,
		"start_time": config.ContestStartTime.Format("2006-01-02 15:04:05"),
		"end_time":   config.ContestEndTime.Format("2006-01-02 15:04:05"),
		"note":       config.ContestNote,
	}
	util.SuccessResponse(w, r, info)
}
