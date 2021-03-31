package main

import (
	"flag"
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/deploy"
	"github.com/HeRaNO/xcpc-team-reg/routers"
)

func main() {
	initRoot := flag.Bool("i", false, "init root user")
	initConfigFile := flag.String("c", "./conf/config.yaml", "the path of configure file")
	flag.Parse()

	config.InitConfig(initConfigFile)

	if *initRoot {
		deploy.InitRootUser()
		return
	}

	err := http.ListenAndServe(config.ListenPort, routers.InitRouters())
	if err != nil {
		panic(err)
	}
}
