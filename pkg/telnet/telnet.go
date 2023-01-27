package telnet

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/kpango/glg"
	"github.com/pkg/errors"
)

func Telnet(ctx context.Context, host string, port int) (retTime float64, err error) {
	now := time.Now()
	defer func() {
		retTime = time.Since(now).Seconds()
		glg.Infof("telnet %s:%s use time %f", host, strconv.Itoa(port), retTime)
		if err != nil {
			glg.Errorf("telnet %s:%s err: %v", host, strconv.Itoa(port), err)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "telnet", host, strconv.Itoa(port))
	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		err = fmt.Errorf("command timed out")
		return
	}

	glg.Infof("telnet %s:%s output: %s", host, strconv.Itoa(port), string(out))
	if strings.Contains(string(out), "Escape character") {
		err = nil
		// contain Escape character is ok, return
		return
	}

	if err != nil {
		err = errors.Wrapf(err, "non-zero exit code")
		return
	}

	return
}
