package player

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func runJXA(ctx context.Context, script string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "osascript", "-l", "JavaScript", "-e", script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		errOut := strings.TrimSpace(stderr.String())
		if strings.Contains(errOut, "-1743") || strings.Contains(strings.ToLower(errOut), "not permitted") {
			return nil, ErrAutomationPermission
		}
		if strings.Contains(errOut, "-600") || strings.Contains(strings.ToLower(errOut), "isn't running") {
			return nil, ErrMusicAppNotRunning
		}
		return nil, fmt.Errorf("jxa execution failed: %w: %s", err, errOut)
	}

	return stdout.Bytes(), nil
}
