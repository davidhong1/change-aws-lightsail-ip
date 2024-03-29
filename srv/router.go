package srv

import (
	"github.com/davidhong1/change-aws-lightsail-ip/config"
	"github.com/davidhong1/change-aws-lightsail-ip/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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
		g.GET("/telnet/:ip", controller.Telnet)
		g.GET("/ip", controller.GetIP)
	}

	// Start server
	e.Logger.Fatal(e.Start(config.Conf.ListenPort))
}
