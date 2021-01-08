package util

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"
)

func ExecutionRestraints(config *Config) {
	// Execution is blocked after this date.
	killDate, err := time.Parse(time.RFC3339, config.Agent.KillDate)
	if err == nil && time.Now().After(killDate) {
		log.Printf("Execution blocked by [Restraint.KillDate] on %s", config.Agent.KillDate)
		os.Exit(1)
	}

	// Execution is blocked when the file is not found in user's home directory.
	if len(config.Restraint.FileKillSwitch) != 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Println("ERROR: [Restraint.FileKillSwitch] unable to resolve user's home directory")
			os.Exit(1)
		}
		if _, err := os.Stat(home + "/" + config.Restraint.FileKillSwitch); err != nil {
			log.Printf("Execution blocked by [Restraint.FileKillSwitch] unable to read file: %s",
				home+"/"+config.Restraint.FileKillSwitch)
			os.Exit(1)
		}
	}

	// Execution is blocked when HTTP/S GET returns 200.
	if len(config.Restraint.HttpKillSwitch) != 0 {
		var data []byte
		client := &http.Client{}
		req, err := http.NewRequest("GET", config.Restraint.HttpKillSwitch, bytes.NewBuffer(data))
		if err != nil {
			log.Printf("ERROR: [Restraint.HttpKillSwitch] unable to build request: %s",
				config.Restraint.HttpKillSwitch)
			os.Exit(1)
		}
		req.Header.Set("User-Agent", config.Agent.Useragent)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR: [Restraint.HttpKillSwitch] failed to parse HTTP response: %s", err)
			os.Exit(1)
		}
		if resp.StatusCode == 200 {
			log.Println("Execution blocked by [Restraint.HttpKillSwitch]")
			os.Exit(1)
		}
	}

	// Execution blocked if the hostname is not listed in config.
	if len(config.Restraint.AllowHost) != 0 {
		allowedHost := false
		hostName, err := os.Hostname()
		if err != nil {
			log.Println("ERROR: [Restraint.AllowHost] unable to check hostname restraint")
			os.Exit(1)
		}
		for _, eachHost := range config.Restraint.AllowHost {
			if eachHost == hostName {
				allowedHost = true
			}
		}
		if !allowedHost {
			log.Println("Execution blocked by [Restraint.AllowHost]")
			os.Exit(1)
		}
	}

	// Execution blocked if the username is not listed in config.
	if len(config.Restraint.AllowUser) != 0 {
		allowedUser := false
		userName, err := user.Current()
		if err != nil {
			log.Println("ERROR: [Restraint.AllowHost] unable to check hostname restraint")
			os.Exit(1)
		}
		for _, eachUser := range config.Restraint.AllowUser {
			if eachUser == userName.Username {
				allowedUser = true
			}
		}
		if !allowedUser {
			log.Println("Execution blocked by [Restraint.AllowedUser]")
			os.Exit(1)
		}
	}
}
