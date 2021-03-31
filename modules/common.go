package modules

import (
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/util"
)

func GetSchool(w http.ResponseWriter, r *http.Request) {
	// return SchoolMap in JSON

	util.SuccessResponseWithTotal(w, r, config.SchoolMap, len(config.SchoolMap))
}

func SayHello(w http.ResponseWriter, r *http.Request) {
	// just say hello

	util.SuccessResponse(w, r, "xcpc-team-reg")
}

func SayHelloAdmin(w http.ResponseWriter, r *http.Request) {
	// just say hello

	util.SuccessResponse(w, r, "xcpc-team-reg admin")
}
