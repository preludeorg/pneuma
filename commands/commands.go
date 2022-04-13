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
func RunCommand(message string, executor string, payloadPath string, agent *util.AgentConfig) (util.Response, util.Process, util.Timeline) {
	switch executor {
	case "keyword":
		task := splitMessage(message, '.')
		switch task[0] {
		default:
			return util.BuildErrorResponse("Keyword selected not available for agent")
		}
	case "config":
		return updateConfiguration(message, agent)
	case "shell":
		return pty.SpawnShell(message, agent)
	case "exit":
		return shutdown(agent)
	default:
		util.DebugLogf("Running instruction")
		return execute(message, executor, agent)
	}
}

func execute(command string, executor string, agent *util.AgentConfig) (util.Response, util.Process, util.Timeline) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(agent.CommandTimeout)*time.Second)
	defer cancel()
	var response util.Response
	var process util.Process
	var timeline util.Timeline
	timeline.Started = time.Now().UnixMilli()
	bites, pid, status := execution(getShellCommand(ctx, executor, command))
	timeline.Finished = time.Now().UnixMilli()
	process.ID = pid
	response.Status = status
	if ctx.Err() == context.DeadlineExceeded {
		bites = []byte("Command timed out.")
	}
	if status == util.SuccessExitStatus {
		response.Output = string(bites)
	} else {
		response.Error = string(bites)
	}
	return response, process, timeline
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

func updateConfiguration(config string, agent *util.AgentConfig) (util.Response, util.Process, util.Timeline) {
	newConfig, err := util.ParseArguments(config)
	if err == nil {
		agent.SetAgentConfig(newConfig)
		return util.BuildStatusResponse("Successfully updated agent configuration.", util.SuccessExitStatus, os.Getpid())
	}
	return util.BuildStatusResponse(err.Error(), 1, os.Getpid())
}

func shutdown(agent *util.AgentConfig) (util.Response, util.Process, util.Timeline) {
	go func(a *util.AgentConfig) {
		time.Sleep(time.Duration(a.KillSleep) * time.Second)
		os.Exit(util.SuccessExitStatus)
	}(agent)
	return util.BuildStatusResponse(fmt.Sprintf("Exiting agent in %d seconds", agent.KillSleep), util.SuccessExitStatus, os.Getpid())
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
