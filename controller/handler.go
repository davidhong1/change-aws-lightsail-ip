package controller

import (
	"fmt"
	"net/http"

	"github.com/davidhong1/change-aws-lightsail-ip/config"
	"github.com/davidhong1/change-aws-lightsail-ip/internal/iplist"
	"github.com/davidhong1/change-aws-lightsail-ip/pkg/telnet"

	"github.com/labstack/echo/v4"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func Telnet(c echo.Context) error {
	ip := c.Param("ip")

	retTime, err := telnet.Telnet(c.Request().Context(), ip, config.Conf.DefaultPort)
	if err != nil {
		return c.String(http.StatusOK, "err: "+err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("use time %fs, %s", retTime, ip))
}

func GetIP(c echo.Context) error {
	ret := ""
	for _, v := range iplist.Get() {
		ret = ret + v + "\n"
	}
	if ret == "" {
		ret = "please try again in 5 minutes"
	}
	return c.String(http.StatusOK, ret)
}
