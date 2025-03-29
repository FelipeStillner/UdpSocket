# Build the application
build:
	go build -o bin/server ./cmd/server/main.go
	go build -o bin/client ./cmd/client/main.go

# Run the server
server:
	./bin/server

# Run the client
client:
	./bin/client

# Clean build files
clean:
	rm -rf bin/
