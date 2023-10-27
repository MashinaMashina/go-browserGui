install-dependencies:
	go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo

build:
	go generate browserGui/cmd/main
	go build -ldflags "-s -H windowsgui" browserGui/cmd/main