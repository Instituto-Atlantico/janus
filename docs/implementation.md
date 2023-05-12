# Implementation docs

All the IoT devices used in janus tests were raspberry pis versions 3 and 4 

## Technologies

Main technologies used in development

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/ubuntu)
- [Docker Compose](https://docs.docker.com/compose/install/linux)

Main technologies used in Raspberry Pi 3/4
- Raspberry Pi OS (64 bit desktop and lite)  
- Docker

Hint to install docker on rasp

```cmd
sudo apt update
sudo apt upgrade
sudo apt install raspberrypi-kernel raspberrypi-kernel-headers
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker pi
sudo reboot
```

## Processes
The two main process of Janus are the Device Provisioning and the Sensor Measurement with presentations. The following are two diagrams showing the implementation of these functionalities:

### Device Provisioning

Device provisioning refers to the deploy of an aca-py agent on the rasp, done after its register on Dojot. To provision a device you need to pass the user@ip of the host and its Dojot id, required for storing sensor measurements in the next step.

Device provisioning refers to the register of an aca-py agent on janus-controller. It's done by an API request /provision with the following information:

```json
{
    "deviceHostName": "pi@192.168.0.1",
    "permissions": ["temperature", "humidity"],
    "brokerIp": "192.168.0.12",
    "brokerUsername": "admin:e72928",
    "brokerPassword": "admin"
}
```

The requirements for this are:
    1. Having docker installed on the device
    2. The device mush already have an aca-py agent deployed, what can be made by janus-cli
    3. The device must be on the same network as the server machine 

``` mermaid
sequenceDiagram
    title: Device provisioning
    autonumber

    participant user
    participant janus as janus-controller
    participant server as server-agent
    participant rasp as rasp-agent

    user ->>+ janus: Ask for device provisioning
    note over user, janus: Permissions and agent and broker information

    janus ->>+ server: Ask for invitation
    server -->>- janus: Done

    janus ->>+ rasp: Send invitation
    rasp -->>- janus: Accepted

    janus ->> janus: Get device permissions

    janus ->>+ server: Ask for credential offer to device
    server ->>+ rasp: Offer credential with permissions
    rasp -->>-server : Accepted
    server -->>-janus: Done

    janus -->>- user: Done
```

### Sensor measurement

The sensor measurement will run periodically, asking for presentation proof to the device and making a validation on it. If the presentation is valid, janus will request the sensor data from the host and send it to Dojot, using the device id passed on device provisioning

``` mermaid
sequenceDiagram
    title: Sensor measurement 
    autonumber

    participant janus as janus-controller
    participant server as server-agent
    participant rasp as rasp-agent

    loop for each x minutes

        janus ->>+ rasp: Collect sensor data
        rasp -->>- janus: returns sensor data
        note over rasp: {"humidity":20.2,"temperature": 32.0, "smoke":false}


        janus ->> janus: get the sensor types received
        loop for each sensor type
            janus ->>+ server: Ask for device presentation-proof with the received sensor types
            server ->> rasp: Request presentation-proof
            rasp ->> rasp: Create presentation proof with credential
            rasp ->> server: Send Presentation Proof
            server -->>- janus: Done
        end

        janus ->> dojot: Send sensor data
    end
```