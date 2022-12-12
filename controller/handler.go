package controller

import (
	"context"
	"fmt"
	"github.com/davidhong1/change-aws-lightsail-ip/config"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/kpango/glg"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func TelnetMe(c echo.Context) error {
	remoteIp := c.RealIP()

	retTime, err := doTelnet(remoteIp)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, "err: "+err.Error())
	}

	// TODO
	config.Conf.DefaultIP = remoteIp

	return c.String(http.StatusOK, fmt.Sprintf("use time %fs, %s", retTime, config.Conf.DefaultIP))
}

func RestIP(c echo.Context) error {
	remoteIP := c.RealIP()
	glg.Info("rest ip %s", remoteIP)
	if remoteIP == config.Conf.DefaultIP {
		config.Conf.DefaultIP = ""
	} else {
		glg.Info("remoteIP not equal to defaultIP")
	}

	return c.String(http.StatusOK, "reset")
}

func Telnet(c echo.Context) error {
	if config.Conf.DefaultIP == "" {
		return c.String(http.StatusOK, "ip is empty")
	}

	retTime, err := doTelnet(config.Conf.DefaultIP)
	if err != nil {
		return c.String(http.StatusOK, "err: "+err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("use time %fs, %s", retTime, config.Conf.DefaultIP))
}

func ChangeIP(c echo.Context) error {
	ip := c.Param("ip")
	config.Conf.DefaultIP = ip

	return c.String(http.StatusOK, "ok")
}

func UpdateIP(c echo.Context) error {
	ip := c.Param("ip")
	config.Conf.DefaultIP = ip

	return c.String(http.StatusOK, "ok")
}

func doTelnet(host string) (retTime float64, err error) {
	now := time.Now()
	defer func() {
		retTime = time.Since(now).Seconds()
		if err != nil {
			glg.Error(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "telnet", host, "10310")
	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		err = fmt.Errorf("command timed out")
		return
	}

	if !strings.Contains(string(out), "Escape character") {
		err = fmt.Errorf("telnet output not contains Escape character")
		return
	}
	glg.Info("Output: %s", string(out))
	if err != nil {
		err = errors.Wrapf(err, "non-zero exit code")
		return
	}

	return
}
