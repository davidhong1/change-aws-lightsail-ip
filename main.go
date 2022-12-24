package main

import (
	"flag"

	"github.com/davidhong1/change-aws-lightsail-ip/config"
	"github.com/davidhong1/change-aws-lightsail-ip/srv"
	"github.com/kpango/glg"
)

var configFilePath = flag.String("config", "", "Must set Config file")

func main() {
	flag.Parse()
	if *configFilePath == "" {
		glg.Fatalln("Must set Config file")
		return
	}

	config.Init(*configFilePath)
	srv.InitTelnetJob()
	srv.Router()
}
