default: build

APP_NAME := "cachex"
CMD_DIR := "cmd/$(APP_NAME)/main.go"

.PHONY: install
install:
	go install github.com/ayuxdev/cachex/cmd/cachex@latest

.PHONY: run
run:
	go run $(CMD_DIR) $(ARGS)

.PHONY: build
build:
	go build -o $(APP_NAME) $(CMD_DIR)

.PHONY: bin
bin:
	mv $(APP_NAME) ~/go/bin/$(APP_NAME)

.PHONY: install-bin
install-bin:
	$(MAKE) build
	$(MAKE) bin