package pty

import (
	"context"
	"github.com/preludeorg/pneuma/util"
	"io"
	"net"
	"runtime"
	"time"
)

func SpawnShell(args string, agent *util.AgentConfig) (string, int, int) {
	var executor string
	switch runtime.GOOS {
	case "windows":
		executor = "powershell.exe"
	case "darwin":
		executor = "/bin/zsh"
	case "linux":
		executor = "/bin/bash"
	default:
		executor = "/bin/sh"
	}
	data := util.ParseArguments(args)
	pid, status, err := spawnPtyShell(data[0], executor, agent)
	if err != nil {
		return err.Error(), status, pid
	}
	return "Successfully spawned shell.", status, pid
}

func cancelOnSocketClose(cancel context.CancelFunc, conn net.Conn) {
	for {
		time.Sleep(time.Duration(30) * time.Second)
		one := make([]byte, 1)
		if _, err := conn.Read(one); err == io.EOF {
			conn.Close()
			cancel()
			return
		}
	}
}