# Define Go command and flags
GO = go
GOFLAGS = -ldflags="-s -w"

TARGET = .bin/go-clean-arch

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
	$(GO) test ./...