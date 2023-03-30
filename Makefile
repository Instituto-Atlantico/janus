WINDOWS=janus-cli_windows_amd64.exe
LINUX=janus-cli_linux_amd64
DARWIN=janus-cli_darwin_amd64

build-cli:
	env GOOS=windows GOARCH=amd64 go build -v -o bin/$(WINDOWS) ./src/janus-cli
	env GOOS=linux GOARCH=amd64 go build -v -o bin/$(LINUX) ./src/janus-cli
	env GOOS=darwin GOARCH=amd64 go build -v -o bin/$(DARWIN) ./src/janus-cli
	@echo build complete