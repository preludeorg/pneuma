package util

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"time"
)

var DebugMode *bool

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
	Name 	  string
	AESKey    []byte
	Range     string
	Contact   string
	Address   string
	Useragent string
	Sleep     int
	KillSleep int
	CommandTimeout int
	Pid int
	Executing map[string]Instruction
}

type Beacon struct {
	Name string
	Target string
	Hostname string
	Location string
	Platform string
	Executors []string
	Range string
	Sleep int
	Pwd string
	Executing string
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
		Name:      pickName(12),
		AESKey:    []byte("abcdefghijklmnopqrstuvwxyz012345"),
		Range:     "red",
		Contact:   "tcp",
		Address:   "127.0.0.1:2323",
		Useragent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36",
		Sleep:     60,
		KillSleep: 5,
		CommandTimeout: 60,
		Pid: os.Getpid(),
		Executing: make(map[string]Instruction),
	}
}

func (c *AgentConfig) SetAgentConfig(ac map[string]interface{}) {
	c.Name = applyKey(c.Name, ac, "Name").(string)
	c.AESKey = applyKey(c.AESKey, ac, "AESKey").([]byte)
	c.Range = applyKey(c.Range, ac, "Range").(string)
	c.Useragent = applyKey(c.Useragent, ac, "Useragent").(string)
	c.Sleep = applyKey(c.Sleep, ac, "Sleep").(int)
	if key, ok := ac["Contact"]; ok {
		if _, ok = CommunicationChannels[strings.ToLower(key.(string))]; ok {
			c.Contact = strings.ToLower(key.(string))
			c.Address = applyKey(c.Address, ac, "Address").(string)
		}
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
	hostname, _ := os.Hostname()
	return Beacon {
		Name:      c.Name,
		Target:	   c.Address,
		Hostname:  hostname,
		Range:     c.Range,
		Sleep:	   c.Sleep,
		Pwd:       pwd,
		Location:  executable,
		Platform:  runtime.GOOS,
		Executors: DetermineExecutors(runtime.GOOS, runtime.GOARCH),
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

func ParseArguments(args string) []string {
	var data []string
	err := json.Unmarshal([]byte(args), &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
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

func applyKey(curr interface{}, ac map[string]interface{}, key string) interface{} {
	if val, ok := ac[key]; ok {
		if key == "Sleep" && reflect.TypeOf(val).Kind() == reflect.Float64 {
			return int(reflect.ValueOf(val).Float())
		}
		if key == "AESKey" && reflect.TypeOf(val).Kind() == reflect.String{
			return []byte(val.(string))
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
