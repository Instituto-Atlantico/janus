# Janus

Janus provides a way to deploy and manage Aries agents on Iot Devices Through a CLI and Aca-py agents.

## Technologies

Main technologies used in Ubuntu desktop

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/ubuntu)
- [Docker Compose](https://docs.docker.com/compose/install/linux)

Main technologies used in Raspberry Pi 3/4
- Raspberry Pi OS (64 bit)
- Docker

Hint with commands to install

```
sudo apt update
sudo apt upgrade
sudo apt install raspberrypi-kernel raspberrypi-kernel-headers
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker <device-user>
sudo reboot
```

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

Is possible run the CLI in other operating systems? Yes, it is!!

To run in Windows check [here](./docs/windows.md) the doc.

To run in Apple check [here](./docs/apple.md) the doc.

Read more about the proposed features [here](./docs/readme.md)
