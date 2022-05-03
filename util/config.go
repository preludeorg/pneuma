package util

import (
	"bytes"
	"crypto/md5"
	"embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"time"
)

const (
	ErrorExitStatus   = -1
	SuccessExitStatus = 0
)

//go:embed conf/default.json
var defaultConfig embed.FS

var (
	DebugMode *bool
	_         = reflect.TypeOf(AgentConfig{})
	_         = reflect.TypeOf(Beacon{})
)

//CommunicationChannels contains the contact implementations
var CommunicationChannels = map[string]Contact{}

//Contact defines required functions for communicating with the server
type Contact interface {
	Communicate(agent *AgentConfig, beacon Beacon) (Beacon, error)
}

type Configuration interface {
	ApplyConfig(ac map[string]interface{})
	BuildBeacon() Beacon
}

type Operation interface {
	StartInstructions(instructions []Instruction) (ret []Instruction)
	StartInstruction(instruction Instruction) bool
	EndInstruction(instruction Instruction)
	BuildExecutingHash() string
}

type AgentConfig struct {
	Name           string
	AESKey         string
	Range          string
	Contact        string
	Address        string
	Useragent      string
	Sleep          int
	KillSleep      int
	CommandJitter  int
	CommandTimeout int
	Pid            int
	Proxy          string
	Debug          bool
	Executing      map[string]Instruction
	Executors      []string
}

type Beacon struct {
	Name      string
	Target    string
	Hostname  string
	Location  string
	Platform  string
	Executors []string
	Range     string
	Sleep     int
	Pwd       string
	Executing string
	Links     []Instruction
}

type Instruction struct {
	ID       string `json:"ID"`
	Executor string `json:"Executor"`
	Payload  string `json:"Payload"`
	Request  string `json:"Request"`
	Response string
	Status   int
	Pid      int
}

func BuildAgentConfig() *AgentConfig {
	var agent AgentConfig
	data, _ := defaultConfig.ReadFile("conf/default.json")
	json.Unmarshal(data, &agent)
	agent.Name = pickName(12)
	agent.Pid = os.Getpid()
	agent.Executors = DetermineExecutors(runtime.GOOS, runtime.GOARCH)
	agent.Executing = make(map[string]Instruction)
	return &agent
}

func (c *AgentConfig) SetAgentConfig(ac map[string]interface{}) {
	c.Name = ApplyKey(c.Name, ac, "Name").(string)
	c.AESKey = ApplyKey(c.AESKey, ac, "AESKey").(string)
	c.Range = ApplyKey(c.Range, ac, "Range").(string)
	c.Useragent = ApplyKey(c.Useragent, ac, "Useragent").(string)
	c.Proxy = ApplyKey(c.Proxy, ac, "Proxy").(string)
	c.Sleep = ApplyKey(c.Sleep, ac, "Sleep").(int)
	c.CommandJitter = ApplyKey(c.CommandJitter, ac, "CommandJitter").(int)
	c.CommandTimeout = ApplyKey(c.CommandTimeout, ac, "CommandTimeout").(int)
	c.Contact = ApplyKey(c.Contact, ac, "Contact").(string)
	c.Address = ApplyKey(c.Address, ac, "Address").(string)
	if _, ok := ac["RefreshExecutors"]; ok {
		c.Executors = DetermineExecutors(runtime.GOOS, runtime.GOARCH)
	}
}

func (c *AgentConfig) StartInstruction(instruction Instruction) bool {
	if _, ex := c.Executing[instruction.ID]; ex {
		return false
	}
	c.Executing[instruction.ID] = instruction
	return true
}

func (c *AgentConfig) StartInstructions(instructions []Instruction) (ret []Instruction) {
	for _, i := range instructions {
		if c.StartInstruction(i) {
			ret = append(ret, i)
		}
	}
	return
}

func (c *AgentConfig) EndInstruction(instruction Instruction) {
	delete(c.Executing, instruction.ID)
}

func (c *AgentConfig) BuildExecutingHash() string {
	if count := len(c.Executing); count > 0 {
		ids := make([]string, count)
		for id := range c.Executing {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		h := md5.New()
		for _, s := range ids {
			io.WriteString(h, s)
		}
		return hex.EncodeToString(h.Sum(nil))
	}
	return ""
}

func (c *AgentConfig) BuildBeacon() Beacon {
	pwd, _ := os.Getwd()
	executable, _ := os.Executable()
	_ = os.Chdir(filepath.Dir(executable))
	hostname, _ := os.Hostname()
	return Beacon{
		Name:      c.Name,
		Target:    c.Address,
		Hostname:  hostname,
		Range:     c.Range,
		Sleep:     c.Sleep,
		Pwd:       pwd,
		Location:  executable,
		Platform:  runtime.GOOS,
		Executors: c.Executors,
		Executing: "",
		Links:     make([]Instruction, 0),
	}
}

func (c *AgentConfig) BuildSocketBeacon(shell string) ([]byte, error) {
	magic := []byte(".p.s.\\")
	header, err := json.Marshal(map[string]string{"name": c.Name, "shell": shell})
	if err != nil {
		return nil, err
	}
	size := new(bytes.Buffer)
	if err = binary.Write(size, binary.LittleEndian, int32(len(header))); err != nil {
		return nil, err
	}
	return bytes.Join([][]byte{magic, size.Bytes(), header}, []byte{}), nil
}

func ParseArguments(args string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(args), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DebugLogf(format string, v ...interface{}) {
	if *DebugMode {
		log.Printf(format, v...)
	}
}

func DebugLog(v ...interface{}) {
	if *DebugMode {
		log.Print(v...)
	}
}

func ApplyKey(curr interface{}, ac map[string]interface{}, key string) interface{} {
	if val, ok := ac[key]; ok {
		if key == "Sleep" && reflect.TypeOf(val).Kind() == reflect.Float64 {
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
