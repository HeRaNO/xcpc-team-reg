package main

import (
	"flag"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	initConfigFile := flag.String("c", "./configs/conf.yaml", "the path of configure file")
	flag.Parse()
	srvAddr := internal.InitConfig(initConfigFile)

	h := server.Default(server.WithHostPorts(srvAddr))
	registerRouter(h)
	h.Spin()
}
