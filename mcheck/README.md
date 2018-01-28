<h3>MongoDB Load Test</h3>

A simple load test tool written in GO.  Execute the binary without installing application (such as Java, Python, and/or other 3pp software).

[TOC]

### 1. Download and Execute
The script has pre-built for Linux, MacOS, and Windows environments.

#### 1.1. On Linux
- Download [mcheck-linux-x64](https://github.com/simagix/mongodb-utils/raw/master/mcheck/build/mcheck-linux-x64)
- Run as, for example

  ```
  mcheck-linux-x64 mongodb://localhost:27017
  ```

#### 1.2. On OS X
- Download [mcheck-osx-x64](https://github.com/simagix/mongodb-utils/raw/master/mcheck/build/mcheck-osx-x64)
- Run as, for example

  ```
  mcheck-osx-x64 mongodb://localhost:27017
  ```

#### 1.3. On Windows
- Download [mcheck-win-x64.exe](https://github.com/simagix/mongodb-utils/raw/master/mcheck/build/mcheck-win-x64.exe)
- Run as, for example

  ```
  mcheck-win-x64 mongodb://localhost:27017
  ```

#### 1.4.
Ctrl-C to stop and to perform cleanup functions

### 2. Develop, Build, and Execute

```
go run mcheck.go
```

#### 2.1. Usage
There are a couple of usages, and they are
* Execute a load test
* See a database for demo

```
mcheck-osx-x64 -h

Usage of build/mcheck-osx-x64:
  -batch int
    	ops per batch (default 512)
  -mongoURI string
    	MongoDB URI (default "mongodb://localhost")
  -seed
    	seed a database for demo
  -size int
    	document size (default 4096)
  -t int
    	number of threads (default 1)
```

For example, to run a load test, execute command as

```
mcheck-osx-x64 -mongoURI=mongodb://user:pwd@mydatabase:27017?authSource=admin -t=4 -size=4096 -batch=300
```

To seed data for demo purpose, add `-seed` to the command line

```
mcheck-osx-x64 -seed=true -mongoURI=mongodb://user:pwd@mydatabase:27017?authSource=admin -t=4 -size=4096 -batch=300
```

#### 2.2. Cross Platforms Compile
Run the *build.sh* script to create binaries for Linux, Windows, and OS X.  The *build.sh* script has the following builds.

```
mkdir build
env GOOS=linux GOARCH=amd64 go build -o build/mcheck-linux-x64 mcheck.go
env GOOS=windows GOARCH=amd64 go build -o build/mcheck-win-x64 mcheck.go
env GOOS=darwin GOARCH=amd64 go build -o build/mcheck-osx-x64 mcheck.go
```
