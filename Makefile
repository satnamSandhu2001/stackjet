build:
	go build -o bin/stackjet main.go

install:
	go build -o /usr/local/bin/stackjet main.go

release:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/stackjet-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o dist/stackjet-linux-arm64 main.go

clean:
	rm -rf bin dist
