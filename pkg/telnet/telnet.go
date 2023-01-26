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
		if err != nil {
			glg.Error(err)
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
