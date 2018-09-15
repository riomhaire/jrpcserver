.DEFAULT_GOAL := everything

dependencies:
	@echo Downloading Dependencies
	@go get ./...

build: dependencies
	@echo Compiling Apps
	@echo   --- github.com/riomhaire/jrpcserver 
	@go build github.com/riomhaire/jrpcserver/infrastructure/application/riomhaire/jrpcserver
	@go install github.com/riomhaire/jrpcserver/infrastructure/application/riomhaire/jrpcserver
	@echo Done Compiling Apps

test:
	@echo Running Unit Tests
	@go test ./...

clean:
	@echo Cleaning
	@go clean
	@find . -name "debug.test" -exec rm -f {} \;

everything: clean build test
	@echo Done

devrun: 
	@cd infrastructure/application/riomhaire/jrpcserver
	@go run main.go

