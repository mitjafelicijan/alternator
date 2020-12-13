# sudo apt install -y entr siege

default:
	@echo "Helpers. Check Makefile for more instructions."

watch:
	find -type f | egrep -i "*.go|*.ini" | entr -r go run *.go --build

requirements:
	go get -u -v -f all

benchmark:
	siege -t 10S -i -c 50 http://localhost:8080

build: clean-build amd64 arm
	@echo "Building amd64 and arm version"

clean-build:
	packr clean
	- rm release -Rf

amd64:
	mkdir -p release/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr build -v -a
	mv alternator release/linux-amd64/alternator

arm:
	mkdir -p release/linux-arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 packr build -v -a
	mv alternator release/linux-arm/alternator
