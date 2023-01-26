# janus

## How to run 

To run the demos you will use [docker-compose](https://docs.docker.com/compose/gettingstarted/).

You need to specify what env configuration file are you using and which docker-compose do you want to use (ledger or ledgerless):

``` shell
docker compose --env-file ./demos/.env.holder  -f ./docker-compose.ledgerless.yml -p testingHolder up
```

The flag -p refers to the project name and you must pass a unique value to it if you want to run two or more agents at the same time