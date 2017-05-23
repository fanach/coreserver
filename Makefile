default: build

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
