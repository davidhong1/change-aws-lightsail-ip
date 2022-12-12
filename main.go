package main

import (
	"flag"

	"github.com/davidhong1/change-aws-lightsail-ip/config"
	"github.com/davidhong1/change-aws-lightsail-ip/controller"
	"github.com/kpango/glg"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var configFilePath = flag.String("config", "", "Must set Config file")

func main() {
	flag.Parse()
	if *configFilePath == "" {
		glg.Fatalln("Must set Config file")
		return
	}

	Init(*configFilePath)
	Router()
}

func Init(configFilePath string) {
	loader := aconfig.LoaderFor(&config.Conf, aconfig.Config{
		SkipEnv:   true,
		SkipFlags: true,
		Files:     []string{configFilePath},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})
	if err := loader.Load(); err != nil {
		glg.Fatalln(err)
		return
	}
	config.Conf.Print()

	controller.InitTelnetJob()
}

func Router() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// get ip without proxy
	// https://echo.labstack.com/guide/ip-address/#case-1-with-no-proxy
	e.IPExtractor = echo.ExtractIPDirect()

	g := e.Group("/v1/changeawslightsailipsrv")
	{
		g.GET("/", controller.Hello)
		g.GET("/telnet", controller.Telnet)
		g.POST("/ip", controller.ChangeIP)
		g.DELETE("/ip", controller.RestIP)
		g.PUT("/ip/:ip", controller.UpdateIP)
		g.GET("/telnetme", controller.TelnetMe)
	}

	// Start server
	e.Logger.Fatal(e.Start(config.Conf.ListenPort))
}
