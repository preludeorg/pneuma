package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/preludeorg/pneuma/util"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

//RunCommand executes a given command
func RunCommand(message string, executor string, payloadPath string, agent *util.AgentConfig) (string, int, int) {
	if strings.HasPrefix(message, "cd") {
		pieces := strings.Split(message, "cd")
		bites := changeDirectory(pieces[1])
		return string(bites), 0, 0
	} else if executor == "keyword" {
		task := splitMessage(message, '.')
		if task[0] == "config" {
			return updateConfiguration(task[1], agent)
		} else if task[0] == "exit" {
			return shutdown(agent)
		}
		return "Keyword selected not available for agent", 0, 0
	} else {
		util.DebugLogf("Running instruction")
		bites, status, pid := execute(message, executor, agent)
		return string(bites), status, pid
	}
}

func execute(command string, executor string, agent *util.AgentConfig) ([]byte, int, int) {
	var bites []byte
	var pid int
	var status int
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(agent.CommandTimeout) * time.Second)
	defer cancel()
	command = convertCommand(executor, command)
	if runtime.GOOS == "windows" {
		if executor == "cmd" {
			bites, pid, status = execution(exec.CommandContext(ctx, "cmd.exe", "/c", command))
		} else {
			bites, pid, status = execution(exec.CommandContext(ctx, "powershell.exe", "-ExecutionPolicy", "Bypass", "-C", command))
		}
	} else {
		if executor == "python" {
			bites, pid, status = execution(exec.CommandContext(ctx, "python", "-c", command))
		} else if executor == "osa" && runtime.GOOS == "darwin" {
			bites, pid, status = execution(exec.CommandContext(ctx, "osascript", "-e", command))
		} else {
			bites, pid, status = execution(exec.CommandContext(ctx, "sh", "-c", command))
		}
	}
	if ctx.Err() == context.DeadlineExceeded {
		bites = []byte("Command timed out.")
	}
	return []byte(fmt.Sprintf("%s%s", bites, "\n")), status, pid
}

func changeDirectory(target string) []byte {
	os.Chdir(strings.TrimSpace(target))
	return []byte(" ")
}

func execution(command *exec.Cmd) ([]byte, int, int){
	var bites []byte
	var status int
	var pid int
	command.SysProcAttr = getSysProcAttrs()
	if out, err := command.Output(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			bites = exitError.Stderr
			pid = exitError.Pid()
			status = exitError.ProcessState.ExitCode()
		} else {
			bites = []byte(err.Error())
			pid = -1
			status = command.ProcessState.ExitCode()
		}
	} else {
		bites = out
		pid = command.ProcessState.Pid()
		status = command.ProcessState.ExitCode()
	}
	return bites, pid, status
}

func updateConfiguration(config string, agent *util.AgentConfig) (string, int, int) {
	var newConfig map[string]interface{}
	err := json.Unmarshal([]byte(config), &newConfig)
	if err == nil {
		agent.SetAgentConfig(newConfig)
		return "Successfully updated agent configuration.", 0, os.Getpid()
	}
	return err.Error(), 1, os.Getpid()
}

func shutdown(agent *util.AgentConfig) (string, int, int) {
	go func(a *util.AgentConfig) {
		time.Sleep(time.Duration(a.KillSleep) * time.Second)
		os.Exit(0)
	}(agent)
	return fmt.Sprintf("Exiting agent in %d seconds", agent.KillSleep), 0, os.Getpid()
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

func convertCommand(executor string, command string) string {
	parts := strings.Split(strings.TrimSpace(command), "\n")
	delimiters := map[string]string{
		"psh":    ";",
		"pwsh":   ";",
		"cmd":    "&",
		"sh":     ";",
		"python": ";",
	}
	if delimiter, ok := delimiters[executor]; len(parts) > 1 && ok {
		return splitAndTrim(parts, delimiter)
	}
	return command
}

func splitAndTrim(parts []string, delimiter string) string {
	n := len(delimiter) * (len(parts) - 1)
	for i := 0; i < len(parts); i++ {
		n += len(parts[i])
	}

	var c strings.Builder
	c.Grow(n)
	c.WriteString(strings.TrimSuffix(parts[0], delimiter))
	for _, p := range parts[1:] {
		c.WriteString(delimiter)
		c.WriteString(strings.TrimSuffix(p, delimiter))
	}

	return c.String()
}