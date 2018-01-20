#! /bin/bash
#
# author: Ken Chen
# create certificates with alternative names
#
PWD=$(pwd)
if [ "$1" = "" ]; then
    TMP=$PWD
else
    TMP=$1
fi

mkdir -p $TMP/certs
cd $TMP/certs
echo "Files are created in $TMP/certs"

read -r -d '' DN <<-EOF
C=US
ST=Georgia
L=Atlanta
O=Simagix
OU=DEV
CN=localhost
emailAddress=ken.chen@simagix.com
EOF

read -r -d '' CADATA <<-EOF
[req]
default_bits = 2048
prompt = no
distinguished_name = dn
default_md = x509
x509_extensions = v3_req

[dn]
$DN

[v3_req]
subjectAltName = @alt_names
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer
basicConstraints = CA:TRUE

[alt_names]
DNS.1=localhost
IP.1=127.0.0.1
EOF

read -r -d '' PEMDATA <<-EOF
[req]
default_bits = 2048
prompt = no
distinguished_name = dn
default_md = x509
req_extensions = v3_req

[dn]
$DN

[v3_req]
EOF

read -r -d '' EXTDATA <<-EOF
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1=localhost
DNS.2=$(hostname -f)
EOF

# CA certificates
echo "Creating server certificate and key file: ca.crt and ca.key"
openssl req -nodes -x509 -days 365 -newkey rsa:2048 -keyout ca.key -out ca.crt -config <(
cat <<-EOF
$CADATA
EOF
)

# Server certificates
openssl req -nodes -newkey rsa:2048 -keyout server.key -out server.csr -config <(
cat <<-EOF
$PEMDATA
EOF
)

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 365 -out server.crt -extfile <(
cat <<-EOF
$EXTDATA
EOF
)
echo "Creating server PEM file: cat server.key server.crt > server.pem"
cat server.key server.crt > server.pem

# Client certificates
echo "Creating client certificate and key file: client.crt and client.key"
openssl req -nodes -newkey rsa:2048 -keyout client.key -out client.csr -config <(
cat <<-EOF
$PEMDATA
EOF
)
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAserial ca.srl -days 365 -out client.crt -extfile <(
cat <<-EOF
$EXTDATA
EOF
)
echo "Creating client PEM file: cat client.key client.crt > client.pem"
cat client.key client.crt > client.pem

ls -l $TMP/certs/ca.crt $TMP/certs/server.pem $TMP/certs/client.pem
cd $PWD
