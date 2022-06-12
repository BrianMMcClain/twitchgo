.DEFAULT_GOAL := build

#help: @ List available tasks on this project
help:
	@echo "user tasks:"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\t\033[32m%-30s\033[0m %s\n", $$1, $$2}'
	@echo

#build: @ Build the twitchgo client binary
build:
	go build

#test: @ Run all tests
test:
	go test ./twitch -race -covermode=atomic -v

#clean: @ Remove build artifacts
clean:
	rm twitchgo