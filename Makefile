.PHONY: build install uninstall clean lint test

GO=go
GOFLAGS=-ldflags="-s -w"
BINARY_NAME=browsercli

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

lint:
	golangci-lint run ./...

test:
	$(GO) build ./...

install: build
	sudo install -m 755 $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

uninstall:
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)