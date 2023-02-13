build-cli: 
	go generate ./src/janus-cli 
	go build ./src/janus-cli 
