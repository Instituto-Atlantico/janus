# janus

Janus provides a way to deploy and manage Aries agents on Iot Devices Through a CLI and Aca-py agents.

## How to build the CLI

Using janus directly through go run is not recommended because some embed and configurations are made on build process. 

To build run:

``` 
make build-cli
```

Alternatively you can run the following: 

```
go generate ./...
go build ./src/janus-cli
```

## How to deploy an agent on a remote device

Having an ssh key par configured and already passed as authorized_keys on remote device is required. Need help with this? Click [here](https://phoenixnap.com/kb/ssh-with-key).

```
./janus-cli deploy remote --agent-port 8001 --agent-name demo -H user@192.168.0.2
```

the aries agent will be available at port 8001 and the admin page at 8002

## How to deploy an agent locally

To deploy locally you can run, but it will only be able to communicate with other local agents.

```
./janus-cli deploy local --agent-port 8001 --agent-name demo
```

If you want to have a communication between local and remote devices you need to pass the network ip for local device:

```
./janus-cli deploy local --agent-port 8001 --agent-name demo --agent-ip 192.168.0.1
```

## Features docs

Read more about the proposed features [here](./docs/readme.md)