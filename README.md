# Pneuma

Pneuma is a cross-compiled Prelude agent which can communicate over all available
Operator listening post protocols, currently TCP, UDP and HTTP.

## Getting started

Clone this repository. Then ensure GoLang 1.13+ is installed before continuing.

To use the agent, install GoLang then start the agent against whichever protocol you want:
```
go run main.go -contact tcp -address 127.0.0.1:2323
go run main.go -contact udp -address 127.0.0.1:4545
go run main.go -contact http -address http://127.0.0.1:3391
```

Change the address field to the location of the Prelude Operator, if you are running your agent on a different computer.

### Note on UDP

Because UDP is a stateless protocol, beacons will be 1-way only, meaning you cannot use this protocol for a full 
adversary emulation exercise. It currently cannot receive instructions. This listening post is designed to be a 
backup <i>heartbeat</i> to prove you still have access to a given computer.

## Compile

When you are ready to use Pneuma in a real environment, you will want to compile it into a binary by running the build.sh script, passing in any string as your unique (public) key, which ensures each compiled agent gets a different file hash:
```
./build.sh JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA
```
GOOS represents the target platform. It can be either darwin, linux or windows.

This will output a file (into the payloads directory) for each supported operating system, which you can copy to any target system and execute normally
to start the agent. 

> Before you compile, consider changing the encryptionKey variable inside cryptic.go. This value represents
> the encryption key to encrypt/decrypt communications with Prelude Operator. This key must be 32-characters
> and must match the encryption key in the Prelude Operator Settings -> local settings section. Also consider
> changing the default address parameters in main.go, so you can start your agent without command-line arguments.

## Use without Operator

While Pneuma is designed to work with the Prelude Operator, as an open-source agent you can point it against any command-and-control (C2) listening post you want. To do this, follow these instructions:

- Ensure your C2 is up & accepting traffic on the same Pneuma port you want to use
- Start Pneuma, pointing it at your C2
- Your C2 will receive an encrypted string, using the KEY inside cryptic.go. Your C2 should use the same encryption key to decrypt the string, which will resolve in a JSON-object, which will contain the below structure:
```
{
  "Name": "test",
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
