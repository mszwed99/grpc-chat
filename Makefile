BINARY_NAME=GrpcChat.app
APP_NAME=gRRC-Chat
VERSION=1.0.0



## Generate Protobuf files
protos gen:
	./chat-protos/generate.sh

## build: build binary and package app
build client:
	rm -rf ${BINARY_NAME}
	rm -f  gRPC-Chat/client
	fyne package -appVersion ${VERSION} -name ${APP_NAME} -release

## run: build and runs the application
run client:
	go run client/.

## clean: runs go clean and deletes binaries
clean client:
	@echo "Cleaning Client..."
	@go clean client
	@rm -rf client/${BINARY_NAME}
	@echo "Cleaned!"

test client:
	go test -v ./client/...

## ------------------------------------------------------------------------
build server:
	rm -rf ${BINARY_NAME}
	rm -f  gRPC-Chat/client

## run: build and runs the application
run server:
	go run server/server.go

## clean: runs go clean and deletes binaries
clean server:
	@echo "Cleaning Server..."
	@go clean server
	@rm -rf server/${BINARY_NAME}
	@echo "Cleaned!"

test client:
	go test -v ./server/...
