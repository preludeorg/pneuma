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
func RunCommand(message string, executor string) (string, int, int) {
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
	var err error
	var status int
	if runtime.GOOS == "windows" {
	    if executor == "cmd" {
			bites, pid, err = execution(exec.Command("cmd.exe", "/c", command))
	    } else {
			bites, pid, err = execution(exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-C", command))
	    }
	} else {
		if executor == "python" {
			bites, pid, err = execution(exec.Command("python", "-c", command))
		} else {
			bites, pid, err = execution(exec.Command("sh", "-c", command))
		}
    }
    if err != nil {
	   bites = []byte(err.Error())
	   status = 1
	}
	return []byte(fmt.Sprintf("%s%s", bites, "\n")), status, pid
}

func changeDirectory(target string) []byte {
	os.Chdir(strings.TrimSpace(target))
	return []byte(" ")
}

func execution(command *exec.Cmd) ([]byte, int, error){
	bites, err := command.Output()
	return bites, command.Process.Pid, err
}