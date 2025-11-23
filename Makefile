.PHONY: build run build-linux clean

APP_NAME=provbot
CMD_PATH=cmd/bot/main.go
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux $(CMD_PATH)

clean:
	rm -rf $(BUILD_DIR)
