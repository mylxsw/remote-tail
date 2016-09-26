
export GOPATH:=$(shell pwd):$(GOPATH)

run:
	go run src/main.go -conf=example.toml

mac:
	go build -o bin/remote-tail-mac src/main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/remote-tail-linux	src/main.go

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/remote-tail-win.exe src/main.go

deploy:
	scp ./bin/remote-tail-linux root@192.168.1.100:/usr/bin/remote-tail
	scp ./bin/remote-tail-linux root@192.168.1.237:/usr/bin/remote-tail

deploy-local:
	cp ./bin/remote-tail-mac /usr/local/bin/remote-tail

clean:
	rm -fr ./bin/remote-tail-linux ./bin/remote-tail-mac ./bin/remote-tail-win.exe
