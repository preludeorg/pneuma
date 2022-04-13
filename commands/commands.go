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
	response, process, timeline := execution(ctx, executor, command)
	if ctx.Err() == context.DeadlineExceeded {
		return util.BuildErrorResponse("Command timed out.")
	}
	return response, process, timeline
}

func execution(ctx context.Context, executor, cmd string) (util.Response, util.Process, util.Timeline) {
	var response util.Response
	var process util.Process
	var timeline util.Timeline
	timeline.Started = time.Now().UnixMilli()
	command := getShellCommand(ctx, executor, cmd)
	out, err := command.Output()
	timeline.Finished = time.Now().UnixMilli()
	response.Output = string(out)
	response.Status = command.ProcessState.ExitCode()
	process.ID = command.ProcessState.Pid()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			response.Error = string(exitError.Stderr)
			response.Status = exitError.ProcessState.ExitCode()
			process.ID = exitError.Pid()
			return response, process, timeline
		}
		response.Error = string(err.Error())
		return response, process, timeline
	}
	return response, process, timeline
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
