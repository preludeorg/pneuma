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
	"time"
)

var UA *string

type HTTP struct{}

func init() {
	util.CommunicationChannels["http"] = HTTP{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func (contact HTTP) Communicate(agent *util.AgentConfig, beacon util.Beacon) (util.Beacon, error) {
	if _, err := checkValidHTTPTarget(agent.Address); err != nil {
		return beacon, err
	}

	setHTTPProxyConfiguration(agent)

	for {
		refreshBeacon(agent, &beacon)
		for agent.Contact == "http" {
			body := beaconPOST(agent.Address, beacon)
			beacon.Links = beacon.Links[:0]
			var tempB util.Beacon
			if err := json.Unmarshal(body, &tempB); err != nil || len(tempB.Links) == 0 {
				break
			}
			runLinks(&tempB, &beacon, agent, "")
		}
		if agent.Contact != "http" {
			return beacon, nil
		}
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

func setHTTPProxyConfiguration(agent *util.AgentConfig) {
	var proxyUrlFunc func(*http.Request) (*url.URL, error)

	if proxyUrl, err := url.Parse(agent.Proxy); err == nil && proxyUrl.Scheme != "" && proxyUrl.Host != "" {
		proxyUrlFunc = http.ProxyURL(proxyUrl)
	} else {
		proxyUrlFunc = http.ProxyFromEnvironment
	}

	http.DefaultTransport.(*http.Transport).Proxy = proxyUrlFunc
}

func requestHTTPPayload(address string, filehash string) ([]byte, int, error) {
	if valid, err := checkValidHTTPTarget(address); !valid {
		return nil, 0, err
	}
	body, _, code, netErr := requestWithHeaders(address, "GET", []byte{}, 1800, map[string]string{"existing": filehash})
	return body, code, netErr
}

func beaconPOST(address string, beacon util.Beacon) []byte {
	data, _ := json.Marshal(beacon)
	body, _, code, err := request(address, "POST", util.Encrypt(data), 20)
	if len(body) > 0 && code == 200 && err == nil {
		return []byte(util.Decrypt(string(body)))
	}
	return body
}

func request(address string, method string, data []byte, timeout time.Duration) ([]byte, http.Header, int, error) {
	return requestWithHeaders(address, method, data, timeout, map[string]string{})
}

func requestWithHeaders(address string, method string, data []byte, timeout time.Duration, headers map[string]string) ([]byte, http.Header, int, error) {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}
	req, err := http.NewRequest(method, address, bytes.NewBuffer(data))
	if err != nil {
		util.DebugLog(err)
	}
	req.Close = true
	req.Header.Set("User-Agent", *UA)
	for header, value := range headers {
		req.Header.Set(header, value)
	}
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
