#! /bin/bash
mkdir -p build
env GOOS=linux GOARCH=amd64 go build -o build/mcheck-linux-x64 mcheck.go
env GOOS=windows GOARCH=amd64 go build -o build/mcheck-win-x64.exe mcheck.go
env GOOS=darwin GOARCH=amd64 go build -o build/mcheck-osx-x64 mcheck.go
