package commands

import (
	"context"
	"fmt"
	"github.com/preludeorg/pneuma/commands/pty"
	"github.com/preludeorg/pneuma/util"
	"os"
	"os/exec"
	"strings"
	"time"
)

//RunCommand executes a given command
func RunCommand(message string, executor string, payloadPath string, agent *util.AgentConfig) (string, int, int) {
	switch executor {
	case "keyword":
		task := splitMessage(message, '.')
		switch task[0] {
		default:
			return "Keyword selected not available for agent", util.ErrorExitStatus, util.ErrorExitStatus
		}
	case "config":
		return updateConfiguration(message, agent)
	case "shell":
		return pty.SpawnShell(message, agent)
	case "exit":
		return shutdown(agent)
	default:
		util.DebugLogf("Running instruction")
		bites, status, pid := execute(message, executor, agent)
		return string(bites), status, pid
	}
}

func execute(command string, executor string, agent *util.AgentConfig) ([]byte, int, int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(agent.CommandTimeout)*time.Second)
	defer cancel()
	bites, pid, status := execution(getShellCommand(ctx, executor, command))
	if ctx.Err() == context.DeadlineExceeded {
		bites = []byte("Command timed out.")
	}
	return []byte(fmt.Sprintf("%s%s", bites, "\n")), status, pid
}

func execution(command *exec.Cmd) ([]byte, int, int) {
	out, err := command.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return append(out, exitError.Stderr...), exitError.Pid(), exitError.ProcessState.ExitCode()
		}
		return append(out, []byte(err.Error())...), util.ErrorExitStatus, command.ProcessState.ExitCode()
	}
	return out, command.ProcessState.Pid(), command.ProcessState.ExitCode()
}

func updateConfiguration(config string, agent *util.AgentConfig) (string, int, int) {
	newConfig, err := util.ParseArguments(config)
	if err == nil {
		agent.SetAgentConfig(newConfig)
		return "Successfully updated agent configuration.", util.SuccessExitStatus, os.Getpid()
	}
	return err.Error(), 1, os.Getpid()
}

func shutdown(agent *util.AgentConfig) (string, int, int) {
	go func(a *util.AgentConfig) {
		time.Sleep(time.Duration(a.KillSleep) * time.Second)
		os.Exit(util.SuccessExitStatus)
	}(agent)
	return fmt.Sprintf("Exiting agent in %d seconds", agent.KillSleep), util.SuccessExitStatus, os.Getpid()
}

func splitMessage(message string, splitRune rune) []string {
	quoted := false
	values := strings.FieldsFunc(message, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == splitRune
	})
	return values
}
