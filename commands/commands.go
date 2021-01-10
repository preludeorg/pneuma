package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//RunCommand executes a given command
func RunCommand(message string, executor string, payloadPath string) (string, int, int) {
	if strings.HasPrefix(message, "cd") {
		pieces := strings.Split(message, "cd")
		bites := changeDirectory(pieces[1])
		return string(bites), 0, 0
	} else if executor == "keyword" {
		return "no keyword executors have been configured in this agent", 0, 0
	} else {
		log.Print("Running instruction")
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

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}