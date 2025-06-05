APP_NAME = stackjet
BIN_DIR = bin
DIST_DIR = dist

build:
	go build -o $(BIN_DIR)/$(APP_NAME) main.go

install:
	go build -o /usr/local/bin/$(APP_NAME) main.go

release:
	mkdir -p $(DIST_DIR)
	GOOS=linux   GOARCH=amd64 go build -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 main.go
	GOOS=linux   GOARCH=arm64 go build -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 main.go

clean:
	rm -rf $(BIN_DIR) $(DIST_DIR)
