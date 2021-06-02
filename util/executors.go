package util

import (
	"fmt"
	"os/exec"
	"strings"
)

type ListFlags []string

func (l *ListFlags) String() string {
	return fmt.Sprint(*l)
}

func (l *ListFlags) Set(value string) error {
	for _, item := range strings.Split(value, ",") {
		*l = append(*l, item)
	}
	return nil
}

//DetermineExecutors looks for available execution engines
func DetermineExecutors(platform string, arch string) []string {
	platformExecutors := map[string]map[string][]string {
		"windows": {
			"file": {"pwsh.exe", "powershell.exe", "cmd.exe"},
			"executor": {"pwsh", "psh", "cmd"},
		},
		"linux": {
			"file": {"python3", "pwsh", "sh", "bash"},
			"executor": {"python", "pwsh", "sh", "bash"},
		},
		"darwin": {
			"file": {"python3", "pwsh", "sh", "osascript", "bash"},
			"executor": {"python", "pwsh", "sh", "osa", "bash"},
		},
	}
	var executors []string
	for platformKey, platformValue := range platformExecutors {
		if platform == platformKey {
			for i := range platformValue["file"] {
				if checkIfExecutorAvailable(platformValue["file"][i]) {
					executors = append(executors, platformExecutors[platformKey]["executor"][i])
				}
			}
		}
	}
	executors = append([]string{"keyword"}, executors...)
	return executors
}

func checkIfExecutorAvailable(executor string) bool {
	_, err := exec.LookPath(executor)
	return err == nil
}