# Define Go command and flags
GO = go
GOFLAGS = -ldflags="-s -w"
GOBIN ?= $$(go env GOPATH)/bin

TARGET = .bin/go-clean-arch
TESTARGS=-coverprofile=./cover.out -covermode=atomic -coverpkg=./...

all: clean $(TARGET)

generate:
	$(GO) generate ./...

$(TARGET):generate
	$(GO) build $(GOFLAGS) -o $(TARGET) cmd/app/main.go

clean:
	rm -f $(TARGET)

run: $(TARGET)
	./$(TARGET)

test:
	$(GO) test ./... $(TESTARGS)

install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

check-coverage: install-go-test-coverage test
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml