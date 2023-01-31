# janus

## How to run 

To run the demos you will use [docker-compose](https://docs.docker.com/compose/gettingstarted/).

1. create docker janus-network

```
./network-up.sh
```

2. run issuer

```
docker compose --env-file ./demos/.env.issuer  -f ./docker-compose.ledger.yml -p issuer up
```

3. run holder

```
docker compose --env-file ./demos/.env.holder -f ./docker-compose.ledgerless.yml -p holder up
```


See [demo.md](./demo.md) to steps about the demo 