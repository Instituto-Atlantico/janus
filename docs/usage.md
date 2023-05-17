# Usage

![The same diagram from the home page readme. A two pieces diagram. The first shows an IoT device sending sensor information directly to Dojot MQTT broker, while the second shows Janus issuing credentials and running presentation proof validations with the IoT device, registering DiDs, credentials and verifying presentations with an Indy blockchain and sending the sensor information to Dojot MQTT broker](./diagram.png)

Janus' usage is based on raspberry pi as our IoT devices, Dojot as our Sensors Measurements MQTT brokers and Docker as our main functional requirement. All the steps for running everything will be described here.

## Raspberry setup

### Configuring the OS

Use [Raspberry Pi Imager](https://www.raspberrypi.com/software/) to instal an OS on the raspberry. We used raspberry pi os lite on our tests. You can also use the imager to pre-connect the device on a wifi network .

Check the following tutorial for help: https://raspberrytips.com/raspberry-pi-wifi-setup/#:~:text=of%20the%20time.-,Use%20Raspberry%20Pi%20Imager,-The%20easiest%20way.

### Changing ssh keys with host

An SSH key Authentication between the host and IoT device is required for janus to work. The keys are used to ssh connect without passwords.

See https://www.digitalocean.com/community/tutorials/how-to-set-up-ssh-keys-2

### Installing docker

Docker is the main functional requirement of janus. To install it on a Raspberry pi device run the following:

```cmd
sudo apt update
sudo apt upgrade
sudo apt install raspberrypi-kernel raspberrypi-kernel-headers
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker pi
sudo reboot
```
### Running a sensor collector 

Here you will need to select what sensors will you send data through Janus. We used to versions of a DHT-11 sensor collector.

A real one (requires a real sensor): https://github.com/eltonlazzarin/rpi-dht11-api-docker
A mocked one (for testing purposes): https://github.com/vitorestevamia/dht11-mock-collector

Will can also, build your ones or use other open source alternatives, but following the requirements:

- The api must use the port 5000 of the raspberry;
- The api must have a single endpoint that returns the value of all sensors you want to send to the broker.

Example:
```cmd
> curl localhost:5000

{"temperature": 21, "humidity": 47}
```

## Dojot

Dojot is the MQTT broker supported by Janus.

### Deploying dojot 

You can check for information about running dojot [here](https://dojotdocs.readthedocs.io/en/latest/installation-guide.html#docker-compose), but in resume:

Clone the docker-compose repo:
```cmd
git clone git@github.com:dojot/docker-compose.git
```

Run it with docker:
```cmd
docker compose --profile complete up --detach
```

Check [this tutorial](https://dojotdocs.readthedocs.io/en/latest/using-web-interface.html#device-management) for details about creating and managing devices with Dojot.

## Janus

One time you had set up the raspberry pi (OS, SSH keys and sensor collectors) and Dojot(Deployment and device creation) you are ready to start working with Janus using janus-cli.

> _**Note:**_ For more details about the CLI use -h flag to get some help:

```cmd
janus-cli -h
janus-cli deploy -h
```

The steps are:

1. Download janus from [github releases page](https://github.com/Instituto-Atlantico/janus/releases) and rename it as janus-cli

```
#linux:
mv ./janus-cli_linux_amd64 ./janus-cli

#mac
mv ./janus-cli_darwin_amd64 ./janus-cli

#windows: 
ren janus-cli_windows_amd64 janus-cli
```

2. Deploy an agent on the rasp:

Use the username and ip of your device as the following <device-username>@<device-ip>

> _**Note:**_ This might take 5-8 minutes to finish.

```cmd
> janus-cli deploy holder -H pi@192.168.0.1
```

3. Deploy the local agent and the janus-controller:

```cmd
> janus-cli deploy issuer 
``` 

4. Make a request to localhost:{controller-port}/provisioning with the device information:

Change the information for matching your device and sensor permissions

> _**Note:**_ You can use localhost:{controller-port}/swagger/ to run the requests in a better ui and for seeing more details

```http
POST http://localhost:{controller-port}/provision HTTP/1.1
content-type: application/json

    {
        "deviceHostName": "pi@192.168.0.1",
        "permissions": ["temperature", "humidity"],
        "brokerIp": "192.168.0.12",
        "brokerUsername": "admin:e72928",
        "brokerPassword": "admin"
    }
```

5. Wait for the device provisioning

6. See the sensor measurements reaching Dojot