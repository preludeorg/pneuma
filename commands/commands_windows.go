package commands

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
)

func getSysProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}

func getShellCommand(ctx context.Context, executor, command string) *exec.Cmd {
	var cmd *exec.Cmd
	switch executor {
	case "python":
		cmd = exec.CommandContext(ctx, "python3.exe", "-c", command)
	case "cmd":
		cmd = exec.CommandContext(ctx, "cmd.exe")
		cmd.SysProcAttr = getSysProcAttrs()
		cmd.SysProcAttr.CmdLine = fmt.Sprintf(`cmd.exe /S /C "%s"`, command)
	default:
		cmd = exec.CommandContext(ctx, "powershell.exe", "-ExecutionPolicy", "Bypass", "-C", command)
		cmd.SysProcAttr = getSysProcAttrs()
	}
	return cmd
}
