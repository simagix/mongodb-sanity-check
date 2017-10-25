# MongoDB Sanity Check
A simple sanity check tool written in GO.  Execute the binary without installing application (such as Java, Python, and/or other 3pp software)
## Download and Execute
### On Linux
- Download [mcheck-linux-x64](https://github.com/simagix/mongodb-sanity-check/blob/master/build/mcheck-linux-x64)
- Run as, for example
```
$ mcheck-linux-64 mongodb://localhost:27017
```
### On OS X
- Download [mcheck-osx-x64](https://github.com/simagix/mongodb-sanity-check/blob/master/build/mcheck-osx-x64)
- Run as, for example
```
$ mcheck-osx-64 mongodb://localhost:27017
```
### On Windows
- Download [mcheck-win-x64](https://github.com/simagix/mongodb-sanity-check/blob/master/build/mcheck-win-x64)
- Run as, for example
```
> mcheck-win-64 mongodb://localhost:27017
```
## Develop, Build, and Execute
```
$ go run mcheck.go
$ go build mcheck.go
$ ./mcheck
```

### Usage
```
$ build/mcheck-osx-x64 -h
Usage of build/mcheck-osx-x64:
  -mongoURI string
    	MongoDB URI (default "mongodb://localhost")
  -t int
    	number of threads (default 1)
  -total int
    	total ops in a batch (default 1000)
```

## Cross Compile
Run the *build.sh* script to create binaries for Linux, Windows, and OS X.  The *build.sh* script has the following builds.
```
$ mkdir build
$ env GOOS=linux GOARCH=amd64 go build -o build/mcheck-linux-x64 mcheck.go
$ env GOOS=windows GOARCH=amd64 go build -o build/mcheck-win-x64 mcheck.go
$ env GOOS=darwin GOARCH=amd64 go build -o build/mcheck-osx-x64 mcheck.go
```
