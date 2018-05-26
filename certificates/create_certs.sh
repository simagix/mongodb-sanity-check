#! /bin/bash
#
# author: Ken Chen
# create certificates with alternative names
#
while getopts ":c:" o; do
    case "${o}" in
        c)
            CA=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [ -z "$*" ]; then
    IFS=', ' read -r -a hostnames <<< $(hostname -f)
else
    hostnames=$@
fi

TMP=$(pwd)
mkdir -p $TMP/certs
rm -f $TMP/certs/*.pem

if [ "$CA" != "" ]; then
    cp $CA $TMP/certs/ca.pem || exit
fi

cd $TMP/certs
echo "Files are created in $TMP/certs"

read -r -d '' DN <<-EOF
C=US
ST=Georgia
L=Atlanta
O=Simagix
OU=DEV
CN=localhost
emailAddress=admin@simagix.com
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

read -r -d '' CLIENT_PEMDATA <<-EOF
[req]
default_bits = 2048
prompt = no
distinguished_name = dn
default_md = x509
req_extensions = v3_req
[dn]
C=US
ST=Georgia
L=Atlanta
O=Simagix
OU=Consulting
CN=ken.chen
emailAddress=ken.chen@simagix.com
[v3_req]
EOF

# CA certificates
# echo "Creating server certificate and key file: ca.crt and ca.key"
if [ "$CA" == "" ]; then
	openssl req -nodes -x509 -days 365 -newkey rsa:2048 -keyout ca.key -out ca.crt -config <(
	cat <<-EOF
	$CADATA
	EOF
	)
	cat ca.crt ca.key > ca.pem
fi

# Server certificates
for hostname in $hostnames
do
	openssl req -nodes -newkey rsa:2048 -keyout server.key -out server.csr -config <(
	cat <<-EOF
	$PEMDATA
	EOF
	)
	
	openssl x509 -req -in server.csr -CA ca.pem -CAkey ca.pem -CAcreateserial -days 365 -out server.crt -extfile <(
	cat <<-EOF
	basicConstraints = CA:FALSE
	keyUsage = nonRepudiation, digitalSignature, keyEncipherment
	subjectAltName = @alt_names
	[alt_names]
	DNS.1=localhost
	DNS.2=127.0.0.1
	DNS.3=$hostname
	EOF
	)
	cat server.key server.crt > ${hostname}.pem

    if [ "$hostname" == "$(hostname -f)" ]; then
        cp ${hostname}.pem server.pem
    fi
done

# Client certificates
# echo "Creating client certificate and key file: client.crt and client.key"
openssl req -nodes -newkey rsa:2048 -keyout client.key -out client.csr -config <(
cat <<-EOF
$CLIENT_PEMDATA
EOF
)
openssl x509 -req -in client.csr -CA ca.pem -CAkey ca.pem -CAserial ca.srl -days 365 -out client.crt -extfile <(
cat <<-EOF
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = TLS Web Server Authentication, TLS Web Client Authentication
EOF
)
cat client.key client.crt > client.pem

rm -f $(ls | grep -v '.pem$')
cd $PWD
ls -l $TMP/certs/
