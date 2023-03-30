# Janus

Janus provides a way to deploy and manage Aries agents on Iot Devices Through a CLI and Aca-py agents.

## How to clone this repository

```bash
git clone https://github.com/Instituto-Atlantico/janus.git
```

## How to build the CLI

Using Janus directly through go run is not recommended because some embed and configurations are made on build process. 

To build run:

```
cd janus && make build-cli
```

## How to deploy an agent on a remote device via Ubuntu terminal

Having an ssh key par configured and already passed as authorized_keys on remote device is required. Need help with this? Click [here](https://phoenixnap.com/kb/ssh-with-key).

```
./bin/janus-cli_linux_amd64 deploy remote --agent-port <port-number> --agent-name <agent-name> -H <device-user>@<device-ip>
```

For example:

```
./bin/janus-cli_linux_amd64 deploy remote --agent-port 8001 --agent-name demo -H pi@192.168.0.2
```

The aries agent will be available at port 8001 and the admin page at port 8002, such as http://192.168.0.2:8002

## How to deploy an agent locally via Ubuntu terminal

To deploy locally you can run, but it will only be able to communicate with other local agents.

```
./bin/janus-cli_linux_amd64 deploy local --agent-port <port-number> --agent-name <agent-name>
```

For example:

```
./bin/janus-cli_linux_amd64 deploy local --agent-port 8001 --agent-name demo
```

The aries agent will be available at port 8001 and the admin page at port 8002, such as http://localhost:8002

If you want to have a communication between local and remote devices you need to pass the network IP for local device:

```
./bin/janus-cli_linux_amd64 deploy local --agent-port <port-number> --agent-name <agent-name> --agent-ip <agent-ip>
```

## How to use the CLI in Windows

Having an ssh key par configured and already passed as authorized_keys on remote device is required. Need help with this? Click [here](https://phoenixnap.com/kb/ssh-with-key).

To use the CLI in Windows is necessary to copy the bin folder generated on Ubuntu via terminal and paste it in a Windows directory previously created, like for example Janus and access it via PowerShell and/or Command Prompt

#### How to deploy an agent on a remote device via PowerShell

```
.\bin\janus-cli_windows_amd64 deploy remote --agent-port <port-number> --agent-name <agent-name> -H <device-user>@<device-ip>
```

For example:

```
C:\Users\username\Desktop\Janus> .\bin\janus-cli_windows_amd64 deploy remote --agent-port 8001 --agent-name demo -H pi@192.168.0.2
```

#### How to deploy an agent locally via PowerShell

```
.\bin\janus-cli_windows_amd64 deploy local --agent-port <port-number> --agent-name <agent-name>
```

#### How to deploy an agent on a remote device via Command Prompt

```
janus-cli_windows_amd64.exe deploy remote --agent-port <port-number> --agent-name <agent-name> -H <device-user>@<device-ip>
```

For example:

```
C:\Users\username\Desktop\Janus\bin> janus-cli_windows_amd64.exe deploy remote --agent-port 8001 --agent-name demo -H pi@192.168.0.2
```

#### How to deploy an agent locally via Command Prompt

```
janus-cli_windows_amd64.exe deploy local --agent-port <port-number> --agent-name <agent-name>
```

## Features docs

Read more about the proposed features [here](./docs/readme.md)
