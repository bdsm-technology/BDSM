bdsm: *.go
	@go build -o bdsm -i -v -ldflags="-X main.version=$(shell git describe --always --long --tags --dirty)"