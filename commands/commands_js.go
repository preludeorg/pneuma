package commands

import (
	"context"
	"os/exec"
	"syscall"
)

func getSysProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func getShellCommand(ctx context.Context, executor, command string) *exec.Cmd {
	var cmd *exec.Cmd
	switch executor {
	default:
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}
	cmd.SysProcAttr = getSysProcAttrs()
	return cmd
}