FROM golang:1.8

ENV PROJECT $GOPATH/src/github.com/zyfdegh/fanach/coreserver
WORKDIR $PROJECT

COPY . $PROJECT

RUN	go test $(go list ./... | grep -v /vendor/) && \
	go build -o coreserver

CMD ["./coreserver"]
