# sudo apt install -y entr siege

default:
	@echo "Helpers. Check Makefile for more instructions."

watch:
	find -type f | egrep -i "*.go|*.ini" | entr -r go run *.go --build

requirements:
	go get -u -v -f all

benchmark:
	siege -t 10S -i -c 50 http://localhost:8080

build: amd64 arm
	@echo "Building arm and amd version"

amd64:
	mkdir -p release/linux-amd64
	packr clean
	- rm release/staticgen-*

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr build -v -a
	mv staticgen release/linux-amd64/alternator

arm:
	mkdir -p release/linux-arm
	packr clean
	- rm release/staticgen-*

	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 packr build -v -a
	mv staticgen release/linux-arm/alternator
