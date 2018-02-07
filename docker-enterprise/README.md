# Build MongoDB Enterprise 3.6 Docker Container

## To Build
```
docker build . -t simagix/mongo:3.6.2-ent
```

## To run
```
docker run -i -p 27017:27017 -t simagix/mongo:3.6.2-ent mongod --sslFIPSMode --bind_ip_all
```
