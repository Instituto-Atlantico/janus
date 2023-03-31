## How to use the CLI in Apple operating systems

Having an ssh key par configured and already passed as authorized_keys on remote device is required. Need help with this? Click [here](https://phoenixnap.com/kb/ssh-with-key).

To use the CLI in Apple operating systems is necessary to copy the bin folder generated on Ubuntu via terminal and paste it on your system directory previously created, like for example Janus and access it via terminal

## How to deploy an agent on a remote device via terminal

```
./bin/janus-cli_darwin_amd64 deploy remote --agent-port <port-number> --agent-name <agent-name> -H <device-user>@<device-ip>
```

For example:

```
./bin/janus-cli_darwin_amd64 deploy remote --agent-port 8001 --agent-name demo -H pi@192.168.0.2
```

The aries agent will be available at port 8001 and the admin page at port 8002, such as http://192.168.0.2:8002

## How to deploy an agent locally via terminal

To deploy locally you can run, but it will only be able to communicate with other local agents.

```
./bin/janus-cli_darwin_amd64 deploy local --agent-port <port-number> --agent-name <agent-name>
```

For example:

```
./bin/janus-cli_darwin_amd64 deploy local --agent-port 8001 --agent-name demo
```

The aries agent will be available at port 8001 and the admin page at port 8002, such as http://localhost:8002

If you want to have a communication between local and remote devices you need to pass the network IP for local device:

```
./bin/janus-cli_darwin_amd64 deploy local --agent-port <port-number> --agent-name <agent-name> --agent-ip <agent-ip>
```

## Features docs

Read more about the proposed features [here](./docs/readme.md)
