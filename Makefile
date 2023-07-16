#msys32
build:
	export GOROOT=/c/Program\ Files/Go
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build -buildmode=c-shared -o tbot3.dll main.go