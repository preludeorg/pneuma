package commands

import (
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
		bites, status, pid := execute(message, executor)
		return string(bites), status, pid
	}
}

func execute(command string, executor string) ([]byte, int, int) {
	var bites []byte
	var pid int
	var status int
	if runtime.GOOS == "windows" {
		if executor == "cmd" {
			bites, pid, status = execution(exec.Command("cmd.exe", "/c", command))
		} else {
			bites, pid, status = execution(exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-C", command))
		}
	} else {
		if executor == "python" {
			bites, pid, status = execution(exec.Command("python", "-c", command))
		} else if executor == "osa" && runtime.GOOS == "darwin" {
			bites, pid, status = execution(exec.Command("osascript", "-e", command))
		} else {
			bites, pid, status = execution(exec.Command("sh", "-c", command))
		}
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
			status = exitError.ExitCode()
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