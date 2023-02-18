NAME = $(shell basename $(shell pwd))

.PHONY: all
dependencies:
	@echo "Installing dependencies..."
	@go get -v -t -d ./...

test: dependencies
	@echo "Testing $(NAME)..."
	@go test -v -cover ./...
