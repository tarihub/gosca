export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: gosca

fmt:
	go fmt ./...

vet:
	go vet ./...

gosca:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/gosca ./cmd/gosca

clean:
	rm -f ./bin/gosca
