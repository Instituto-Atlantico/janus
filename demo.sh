docker ps -aq | xargs docker rm -f

docker -H ssh://pi@192.168.0.3 ps -aq | xargs docker -H ssh://pi@192.168.0.3 rm -f

./janus-cli deploy local --agent-name issuer --agent-port 8001 --agent-ip 192.168.0.10

./janus-cli deploy remote --agent-name holderrasp --agent-port 8001 -H pi@192.168.0.3