# janus

Janus provides a way to deploy and manage Aries agents on Iot Devices Through a CLI and Aca-py agents.

## How to build

``` 
make build-cli
```

Alternatively you can run the following: 

```
go generate ./...
go build ./src/janus-cli
```

## How to run

```
./janus-cli deploy --port 8001 --name demo-agent
```

the aries agent will be available at port 8001 and the admin page at 8002