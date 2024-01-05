#run this using msys2 mingw32
build:
	CC=gcc CGO_ENABLED=1 GOOS=windows GOARCH=386 go build -buildmode=c-shared -o tbot3.dll main.go