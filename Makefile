# sudo apt install -y entr siege

VERSION = "2.0.0"

default:
	@echo "Helpers. Check Makefile for more instructions."
	echo "$(VERSION)"

watch: embed-assets
	find -type f | egrep -i "*.go|*.ini" | entr -r go run *.go --watch --http

requirements:
	go get -u -v -f github.com/jteeuwen/go-bindata/...
	go get -u -v -f all

benchmark:
	siege -t 10S -i -c 50 http://localhost:8080

build: clean-build embed-assets amd64 arm
	@echo "Building amd64 and arm version"

clean-build:
	- rm release -Rf

embed-assets:
	go-bindata posts/... template/... assets/... config.ini

amd64:
	mkdir -p release/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.ldVersion=$(VERSION)'" -o release/linux-amd64/alternator -v -a *.go

arm:
	mkdir -p release/linux-arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-X 'main.ldVersion=$(VERSION)'"  -o release/linux-arm/alternator -v -a *.go
