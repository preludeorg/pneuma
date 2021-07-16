//+build !windows, !js

package commands

import "syscall"

func getSysProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}