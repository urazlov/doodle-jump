APP_NAME=doodlejump
SRC=main.go
APP_DIR=$(APP_NAME).app
CONTENTS_DIR=$(APP_DIR)/Contents
MACOS_DIR=$(CONTENTS_DIR)/MacOS

build:
	go build -o $(APP_NAME) $(SRC)
	mkdir -p $(MACOS_DIR)
	mv $(APP_NAME) $(MACOS_DIR)/

open:
	open ${APP_DIR}

format:
	gofmt -w main.go
