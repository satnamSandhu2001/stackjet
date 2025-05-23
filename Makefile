APP_NAME=go-stack-cli

build:
	go build -o bin/$(APP_NAME) main.go

install:
	go build -o /usr/local/bin/$(APP_NAME) main.go

release:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/$(APP_NAME)-linux main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/$(APP_NAME)-mac main.go
	GOOS=windows GOARCH=amd64 go build -o dist/$(APP_NAME).exe main.go

clean:
	rm -rf bin dist
