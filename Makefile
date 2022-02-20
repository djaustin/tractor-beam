appName = tractor-beam

compile:
	GOOS=linux GOARCH=386 go build -o dist/$(appName)-$(version)-linux-386 main.go
	GOOS=linux GOARCH=arm64 go build -o dist/$(appName)-$(version)-linux-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o dist/$(appName)-$(version)-linux-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/$(appName)-$(version)-darwin-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/$(appName)-$(version)-darwin-amd64 main.go
	GOOS=windows GOARCH=386 go build -o dist/$(appName)-$(version)-windows-386.exe main.go
	GOOS=windows GOARCH=arm64 go build -o dist/$(appName)-$(version)-windows-arm64.exe main.go
	GOOS=windows GOARCH=amd64 go build -o dist/$(appName)-$(version)-windows-amd64.exe main.go