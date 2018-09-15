bdsm: *.go
	@go build -o bdsm -i -v -ldflags="-X main.version=$(git describe --always --long --dirty)"