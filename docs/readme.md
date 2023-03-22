# Implementation docs

The two main initial parts of Janus are the Device Provisioning and the Sensor Measurement with presentations. The following are two diagrams showing the proposed implementation of these functionalities:

## Device Provisioning

Device provisioning refers to the deploy of an aca-py agent on the rasp, done after its register on Dojot. To provision a device you need to pass the user@ip of the host and its Dojot id, required for storing sensor measurements in the next step.

The requirements for this are:
    1. Having docker installed on the device
    2. Having a ssh key par configured to enable ssh connection without passwords
    3. The device must be on the same network as the server machine 

``` mermaid
sequenceDiagram
    title: Device provisioning
    autonumber

    participant user
    participant dojot
    participant janus as janus-controler
    participant server as server-agent
    participant rasp as rasp-agent

    user ->> dojot: Register device and sensors
    dojot -->> user: Registered

    user ->> janus: Ask for device registration
    note over user, janus: device registration requires the host ip and the its Dojot id

    janus ->>+ rasp: Deploy aca-py agent
    rasp -->>- janus: Aca-py agent up and running

    janus ->>+ server: Ask for invitation
    server -->>- janus: Done

    janus ->>+ rasp: Ask for invitation acception
    rasp -->>- janus: Done

    janus ->> janus: Get device permissions

    janus ->>+ server: Ask for credential offer to device
    server ->>+ rasp: Offer credential with permissions
    rasp -->>-server : Accepted
    server -->>-janus: Done

    janus -->> user: Done
```

## Sensor measurement

The sensor measurement will run periodically, asking for presentation proof to the device and making a validation on it. If the presentation is valid, janus will request the sensor data from the host and send it to Dojot, using the device id passed on device provisioning

``` mermaid
sequenceDiagram
    title: Sensor measurement 
    autonumber

    participant janus as janus-controler
    participant server as server-agent
    participant rasp as rasp-agent

    loop for each x minutes
        janus ->>+ server: Ask for device presentation-proof
        server ->> rasp: Request presentation-proof
        rasp ->> rasp: Create presentation proof with credential
        rasp ->> server: Send Presentation Proof
        server -->>- janus: Done

        janus ->>+ rasp: Collect sensor data
        rasp -->>- janus: returns sensor data

        janus ->> dojot: Send sensor data
    end
```