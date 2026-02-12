APP_NAME=server
CMD_DIR=cmd/server
BIN_DIR=bin

## Run the server 
run:
	cd $(CMD_DIR) && go run .

## Run all tests
test:
	go test ./...

## Format code
fmt:
	go fmt ./...
