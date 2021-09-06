# Pneuma

Pneuma is a cross-compiled [Operator](https://prelude.org) agent which can communicate over all available listening post protocols, currently gRPC, TCP, UDP and HTTP.

> This repo is for those who feel comfortable working with the inner workings of a Remote Access Trojan (RAT). If you simply need a copy of the Pneuma agent, there are precompiled versions available for [MacOS](https://s3.amazonaws.com/operator.payloads.open/payloads/pneuma/pneuma-darwin), [Windows](https://s3.amazonaws.com/operator.payloads.open/payloads/pneuma/pneuma-windows.exe) and [Linux](https://s3.amazonaws.com/operator.payloads.open/payloads/pneuma/pneuma-linux). You will need to mark these as executable after downloading.

This agent can be used with Operator, as designed, or pointed at any command-and-control (C2) listening post that can accept its beacon. More details on this below.

## Getting started

Clone this repository. Then ensure GoLang 1.16+ is installed before continuing.

To use the agent, install GoLang then start the agent against whichever protocol you want:

> If using a precompiled version, replace 'go run main.go' with './pneuma-windows.exe' in the below commands, replacing pneuma.exe with the name of your downloaded file.

```
go run main.go -contact tcp -address 127.0.0.1:2323
go run main.go -contact udp -address 127.0.0.1:4545
go run main.go -contact http -address http://127.0.0.1:3391
go run main.go -contact grpc -address 127.0.0.1:2513
```

Change the address field to the location of Prelude Operator, if you are running your agent on a different computer.

> gRPC relies on protocol buffers to determine the structure of the data between client <--> server (beacons).
This is statically compiled, so if anything changed in the beacon properties the following command will recompile
the buffer:
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative beacon.proto
```

## Compile

When you are ready to use Pneuma in a real environment, you will want to compile it into a binary by running the build.sh script, passing in any string as your unique (public) randomHash string, which ensures each compiled agent gets a different file hash:
```shell
./build.sh JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA
```

On windows you can build with:

```powershell
.\build.ps1 -randomHash "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"
```

There is also support for garble builds (Thanks to the amazing work here: https://github.com/burrowers/garble):

```shell
GO111MODULE=on go get mvdan.cc/garble
go mod tidy
./garble-build.sh JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA
```

```powershell
GO111MODULE=on go get mvdan.cc/garble
go mod tidy
.\garble-build.ps1 -randomHash "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"
```

This will output a file (into the payloads directory) for each supported operating system, which you can copy to any target system and execute normally
to start the agent. 

> Before you compile, consider changing the AESKey variable inside util/conf/default.json. This value represents
> the encryption key to encrypt/decrypt communications with Prelude Operator. This key must be 32-characters
> and must match the encryption key in the Prelude Operator Emulate section. Also consider
> changing the default address parameters in main.go, so you can start your agent without command-line arguments.

## Use without Operator

While Pneuma is designed to work with Prelude Operator, as an open-source agent you can point it against any command-and-control (C2) listening post you want. To do this, follow these instructions:

- Ensure your C2 is up & accepting traffic on the same Pneuma port you want to use
- Start Pneuma, pointing it at your C2
- The C2 server will send instructions in this format which get encrypted with your KEY inside Key Management in Settings:
```
{
  ID: "067e99fb-f88f-49a8-aadc-b5cadf3438d4",
  ttp: "0b726950-11fc-4b15-a7d3-0d6e9cfdbeab",
  tactic: "discovery",
  Executor: "sh",
  Request: "whoami",
  Payload: "https://s3.amazonaws.com/operator.payloads/demo.exe",
}
```
- Your C2 will receive an encrypted string, using the KEY inside cryptic.go. Your C2 should use the same encryption key to decrypt the string, which will resolve in a JSON-object, which will contain the below structure:
```
{
  "Name": "test",
  "Target": "0.0.0.0:2323",
  "Hostname": "ubuntu-host",
  "Location": "/tmp/me.go"
  "Platform": "darwin",
  "Executors": ["sh"],
  "Range": "red",
  "Pwd": "/tmp",
  "Links": []
}
```

- Your C2 should then fill in links (instructions) for the agent to complete. Each link should have the following structure. Response, Status and Pid will be filled in by Pneuma after running the Request:
```
{
  "ID": "123",
  "Executor": "sh",
  "Payload: "",
  "Request": "whoami",
  "Response: "",
  "Status: 0,
  "Pid": 0
}
```
- Encrypt your new JSON object with the same encryption key and send it back to the agent.

This ping/pong of this beacon/JSON-object should continue on whatever periodic basis you desire. 

## CALDERA inspired

As former CALDERA leads, we wrote the MITRE Sandcat and Manx agents. Pneuma is its own thing but it shares characteristics with both of these original agents.

## Interested in contributing?

We strongly support contributors to this project. Please fork this repo and submit pull requests for considerations.

### Write your own agent

We plan on publishing an agent interface in the future, to describe how to build your own agents in any language. In the meantime, this code was written in a simplistic way to be an example for any agent you want to build.
