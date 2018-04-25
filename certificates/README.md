# Start a sharded MongoDB cluster with TLS/SSL Encryption

Usage: `create_certs.sh [dir]`

```
├── README.md
├── certs
│   ├── ca.crt
│   ├── client.pem
│   └── server.pem
└── create_certs.sh
```

## Generate Certificates

- Generate `ca.crt`
- Generate `server.key` and `server.csr`
- Sign `server.csr` with `ca.crt`
- Create `server.pem` from `server.crt` and `server.key`
- Generate `client.key` and `client.csr`
- Sign `client.csr` with `ca.crt`
- Create `client.pem` from `client.crt` and `client.key`

```
$ create_certs.sh $(pwd)

Files are created in /Users/kenchen/tmp/temp/certs
Creating server certificate and key file: ca.crt and ca.key
Generating a 2048 bit RSA private key
.........................+++
...........................................+++
writing new private key to 'ca.key'
-----
Generating a 2048 bit RSA private key
..+++
.............+++
writing new private key to 'server.key'
-----
Signature ok
subject=/C=US/ST=Georgia/L=Atlanta/O=MongoDB/OU=CE/CN=Kens-MacBook-Pro.local/emailAddress=ken.chen@mongodb.com
Getting CA Private Key
Creating server PEM file: cat server.key server.crt > server.pem
Creating client certificate and key file: client.crt and client.key
Generating a 2048 bit RSA private key
........................................................+++
..........+++
writing new private key to 'client.key'
-----
Signature ok
subject=/C=US/ST=Georgia/L=Atlanta/O=MongoDB/OU=CE/CN=Kens-MacBook-Pro.local/emailAddress=ken.chen@mongodb.com
Getting CA Private Key
Creating client PEM file: cat client.key client.crt > client.pem
```

## Start MongoDB sharded cluster
```
$ mlaunch init --replicaset --nodes 3 --sharded 2 --auth --sslPEMKeyFile certs/server.pem --sslMode allowSSL --sslCAFile certs/ca.crt --sslClientPEMKeyFile certs/client.pem --sslClientCertificate certs/client.crt

launching: mongod on port 27018
launching: mongod on port 27019
launching: mongod on port 27020
launching: mongod on port 27021
launching: mongod on port 27022
launching: mongod on port 27023
launching: config server on port 27024
replica set 'configRepl' initialized.
replica set 'shard01' initialized.
replica set 'shard02' initialized.
launching: mongos on port 27017
adding shards. can take up to 30 seconds...
sent signal 15 to 8 processes.
launching: config server on port 27024
launching: mongod on port 27018
launching: mongod on port 27019
launching: mongod on port 27020
launching: mongod on port 27021
launching: mongod on port 27022
launching: mongod on port 27023
launching: mongos on port 27017
Username "user", password "password"
```

## Connect from Client
```
$ mongo mongodb://user:password@localhost/test?authSource=admin --ssl --sslPEMKeyFile certs/client.pem --sslCAFile certs/ca.crt

MongoDB shell version v3.6.2
connecting to: mongodb://localhost/test?authSource=admin
MongoDB server version: 3.6.2
Server has startup warnings:
2018-01-19T15:54:17.004-0500 I CONTROL  [main]
2018-01-19T15:54:17.005-0500 I CONTROL  [main] ** WARNING: This server is bound to localhost.
2018-01-19T15:54:17.006-0500 I CONTROL  [main] **          Remote systems will be unable to connect to this server.
2018-01-19T15:54:17.008-0500 I CONTROL  [main] **          Start the server with --bind_ip <address> to specify which IP
2018-01-19T15:54:17.009-0500 I CONTROL  [main] **          addresses it should serve responses from, or with --bind_ip_all to
2018-01-19T15:54:17.010-0500 I CONTROL  [main] **          bind to all interfaces. If this behavior is desired, start the
2018-01-19T15:54:17.011-0500 I CONTROL  [main] **          server with --bind_ip 127.0.0.1 to disable this warning.
2018-01-19T15:54:17.012-0500 I CONTROL  [main]
mongos> sh.status()
--- Sharding Status ---
  sharding version: {
  	"_id" : 1,
  	"minCompatibleVersion" : 5,
  	"currentVersion" : 6,
  	"clusterId" : ObjectId("5a625adafc0bfad9bd027174")
  }
  shards:
        {  "_id" : "shard01",  "host" : "shard01/localhost:27018,localhost:27019,localhost:27020",  "state" : 1 }
        {  "_id" : "shard02",  "host" : "shard02/localhost:27021,localhost:27022,localhost:27023",  "state" : 1 }
  active mongoses:
        "3.6.2" : 1
  autosplit:
        Currently enabled: yes
  balancer:
        Currently enabled:  yes
        Currently running:  no
        Failed balancer rounds in last 5 attempts:  0
        Migration Results for the last 24 hours:
                No recent migrations
  databases:
        {  "_id" : "config",  "primary" : "config",  "partitioned" : true }

mongos>
```

## X509 Authentication
### Server Cert
```
openssl x509 -in certs/client.pem -noout -text

Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 10459429611566710654 (0x91275bd0b6a2bf7e)
    Signature Algorithm: sha1WithRSAEncryption
        Issuer: C=US, ST=Georgia, L=Atlanta, O=Simagix, OU=DEV, CN=localhost/emailAddress=admin@simagix.com
        Validity
            Not Before: Apr 25 11:30:07 2018 GMT
            Not After : Apr 25 11:30:07 2019 GMT
        Subject: C=US, ST=Georgia, L=Atlanta, O=Simagix, OU=DEV, CN=localhost/emailAddress=admin@simagix.com
...

mongod --dbpath data --logpath data/mongod.log --fork --sslCAFile certs/ca.crt --sslPEMKeyFile certs/server.pem --auth --sslMode requireSSL
```

### Client Cert
```
openssl x509 -in certs/client.pem -noout -text

Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 10459429611566710655 (0x91275bd0b6a2bf7f)
    Signature Algorithm: sha1WithRSAEncryption
        Issuer: C=US, ST=Georgia, L=Atlanta, O=Simagix, OU=DEV, CN=localhost/emailAddress=admin@simagix.com
        Validity
            Not Before: Apr 25 11:30:08 2018 GMT
            Not After : Apr 25 11:30:08 2019 GMT
        Subject: C=US, ST=Georgia, L=Atlanta, O=Simagix, OU=Consulting, CN=ken.chen/emailAddress=ken.chen@simagix.com
...
```

### Create User in `$external`
```
mongo mongodb://user:password@localhost/admin?authSource=admin --ssl --sslPEMKeyFile certs/client.pem --sslCAFile certs/ca.crt

db.getSisterDB("$external").runCommand(
  {
    createUser:"emailAddress=ken.chen@simagix.com,CN=ken.chen,OU=Consulting,O=Simagix,L=Atlanta,ST=Georgia,C=US" , 
    roles: [{role: 'root', db: 'admin' }] 
  }
)
```

### Authenticate using X509
```
mongo --host localhost --sslCAFile certs/ca.crt --ssl --sslPEMKeyFile certs/client.pem

db.getSisterDB("$external").auth( 
  { 
    mechanism: "MONGODB-X509", 
    user:"emailAddress=ken.chen@simagix.com,CN=ken.chen,OU=Consulting,O=Simagix,L=Atlanta,ST=Georgia,C=US" 
  }
)
```

or

```
mongo --host localhost --sslCAFile certs/ca.crt --ssl --sslPEMKeyFile certs/client.pem --authenticationMechanism MONGODB-X509 --authenticationDatabase "\$external" -u "emailAddress=ken.chen@simagix.com,CN=ken.chen,OU=Consulting,O=Simagix,L=Atlanta,ST=Georgia,C=US"
```