package util

import (
	"math/rand"
	"reflect"
	"strings"
	"time"
)

//CommunicationChannels contains the contact implementations
var CommunicationChannels = map[string]Contact{}

//Contact defines required functions for communicating with the server
type Contact interface {
	Communicate(agent *AgentConfig, beacon Beacon) Beacon
}

type Configuration interface {
	ApplyConfig(ac map[string]interface{})
}

type AgentConfig struct {
	Name 	  string
	AESKey    []byte
	Range     string
	Contact   string
	Address   string
	Useragent string
	Sleep     int
}

type Beacon struct {
	Name string
	Location string
	Platform string
	Executors []string
	Range string
	Sleep int
	Pwd string
	Links []Instruction
}

type Instruction struct {
	ID string `json:"ID"`
	Executor string `json:"Executor"`
	Payload string `json:"Payload"`
	Request string `json:"Request"`
	Response string
	Status int
	Pid int
}

func BuildAgentConfig() *AgentConfig {
	return &AgentConfig{
		Name: pickName(12),
		AESKey: []byte("abcdefghijklmnopqrstuvwxyz012345"),
		Range: "red",
		Contact: "tcp",
		Address: "127.0.0.1:2323",
		Useragent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36",
		Sleep: 60,
	}
}

func (c *AgentConfig) SetAgentConfig(ac map[string]interface{}) {
	c.Name = applyKey(c.Name, ac, "Name").(string)
	c.AESKey = applyKey(c.AESKey, ac, "AESKey").([]byte)
	c.Range = applyKey(c.Range, ac, "Range").(string)
	c.Useragent = applyKey(c.Useragent, ac, "Useragent").(string)
	c.Sleep = applyKey(c.Sleep, ac, "Sleep").(int)
	if key, ok := ac["Contact"]; ok {
		if _, ok = CommunicationChannels[key.(string)]; ok {
			c.Contact = strings.ToLower(key.(string))
			c.Address = applyKey(c.Address, ac, "Address").(string)
		}
	}
}

func applyKey(curr interface{}, ac map[string]interface{}, key string) interface{} {
	if val, ok := ac[key]; ok {
		if key == "Sleep" && reflect.TypeOf(ac[key]).Kind() == reflect.Float64 {
			return int(reflect.ValueOf(val).Float())
		}
		return val
	}
	return curr
}

func pickName(chars int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, chars)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}