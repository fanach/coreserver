default: help

help:
	@echo -e "Select a sub command \n"
	@echo -e "dep-init: \n\t Install govendor and init vendor (after you cloned this repo)"
	@echo -e "dep-update: \n\t Refresh packages under vendor (when you changed imports packages)"
	@echo -e "build: \n\t Build coreserver Docker image"
	@echo -e "run: \n\t Run coreserver Docker container"
	@echo -e "push: \n\t Push coreserver Docker image to DockerHub"
	@echo -e "build-local: \n\t Build coreserver binary to bin/"
	@echo -e "run-local: \n\t Execute local coreserver binary"
	@echo -e "test-local: \n\t Run unit testing in local"
	@echo -e "clean-local: \n\t Remove local binaries, configs"
	@echo -e "\n"
	@echo -e "See README.md for more."

# Build, test, run in local

dep-init:
	go get github.com/kardianos/govendor
	govendor init

dep-update:
	govendor remove +unused
	govendor add +external

build-local:
	go build -o bin/coreserver

run-local:
	./bin/coreserver

test-local:
	go test ./...

clean-local:
	rm -rf bin/*

# Build, test, run in Docker container

build:
	docker build -t fanach/coreserver .

run:
	docker run --rm -p 8080:8080 fanach/coreserver

push: build
	docker push fanach/coreserver .
