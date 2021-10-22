Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

run:
	go run *.go -conf=example.toml

mac:
	go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/remote-tail-mac *.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/remote-tail-linux	*.go

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/remote-tail-win.exe *.go

dist:
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/remote-tail-linux *.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/remote-tail-linux-arm64 *.go
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/remote-tail-darwin *.go	
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/remote-tail.exe *.go	

deploy-local:
	cp ./bin/remote-tail-mac /usr/local/bin/remote-tail

clean:
	rm -fr ./bin/remote-tail-linux ./bin/remote-tail-mac ./bin/remote-tail-win.exe
