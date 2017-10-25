# MongoDB Sanity Check

```
$ go run mcheck.go
$ go build mcheck.go
$ ./mcheck
```

## Cross Compile
### Linux
```
$ mkdir build
$ env GOOS=linux GOARCH=amd64 go build -o build/mcheck-linux-x64 mcheck.go
$ env GOOS=windows GOARCH=amd64 go build -o build/mcheck-win-x64 mcheck.go
$ env GOOS=darwin GOARCH=amd64 go build -o build/mcheck-osx-x64 mcheck.go
```
