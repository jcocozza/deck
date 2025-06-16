linux:
	GOOS=linux GOARCH=amd64 go build -o deck-linux

windows:
	GOOS=windows GOARCH=amd64 go build -o deck-windows

darwin-amd:
	GOOS=darwin GOARCH=amd64 go build -o deck-darwin-amd

darwin-arm:
	GOOS=darwin GOARCH=arm65go build -o deck-darwin-arm
