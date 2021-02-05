package util

import "syscall"

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32 	 = syscall.NewLazyDLL("user32.dll")
)

func HideConsole() {
	setWindow(0)
}

func ShowConsole() {
	setWindow(1)
}

func setWindow(visible int) {
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	if getConsoleWindow.Find() != nil {
		return
	}

	showWindow := user32.NewProc("ShowWindow")
	if showWindow.Find() != nil {
		return
	}

	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}

	showWindow.Call(hwnd, uintptr(visible))
}