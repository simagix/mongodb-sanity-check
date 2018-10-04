#! /bin/bash
#
# author: Ken Chen
# create certificates with alternative names
#
usage() {
    echo "Usage: $0 [-c <ca.pem>] [-o <output dir>] [host ...]" 1>&2

cat << EOF

Environment variables
	C   country
	ST  state
	L   local/city
	O   organization/company
	OU_SERVER   organization unit/group - server
	OU_CLIENT   organization unit/group - client
	CN_SERVER   common name
	CN_CLIENT   common name
	EMAIL_ADM   admin email
	EMAIL_CLIENT    client email
EOF
    exit 1
}

CERTS_DIR="./certs"

while getopts ":c:o:" o; do
    case "${o}" in
        c)
            CA=${OPTARG}
            ;;
        o)
            CERTS_DIR=${OPTARG}
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

if [ -d "$CERTS_DIR" ]; then
    echo -n "Directory $CERTS_DIR exists, override [y/N]? "
    read y
    if [ "$y" == "y" ]; then
        rm -f $CERTS_DIR/*
    else
        exit
    fi
else
    mkdir -p $CERTS_DIR
fi

if [ "$CA" != "" ]; then
    cp $CA $CERTS_DIR/ca.pem || exit
fi

OPWD=$(pwd)
cd $CERTS_DIR
echo "Files are created in $CERTS_DIR"

C="${C:-US}"
ST="${ST:-Georgia}"
L="${L:-Atlanta}"
O="${O:-Simagix}"
OU_SERVER="${OU_SERVER:-Root}"
OU_CLIENT="${OU_CLIENT:-Root}"
CN_SERVER="${CN_SERVER:-localhost}"
CN_CLIENT="${CN_CLIENT:-ken.chen}"
EMAIL_ADM="${EMAIL_ADM:-admin@simagix.com}"
EMAIL_CLIENT="${EMAIL_CLIENT:-ken.chen@simagix.com}"

read -r -d '' DN <<-EOF
C=$C
ST=$ST
L=$L
O=$O
OU=$OU_SERVER
CN=$CN_SERVER
emailAddress=$EMAIL_ADM
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
C=$C
ST=$ST
L=$L
O=$O
OU=$OU_CLIENT
CN=$CN_CLIENT
emailAddress=$EMAIL_CLIENT
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
    rm -f ca.key
else
    echo "extracting ca.crt"
    openssl crl2pkcs7 -nocrl -certfile ca.pem | openssl pkcs7 -print_certs -out ca.crt || exit
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

    if [ "${hostname:0:2}" == "*." ]; then
	    cat server.key server.crt > ${hostname:2}.pem
    else
	    cat server.key server.crt > ${hostname}.pem
    fi
    rm -f server.csr server.crt server.key

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
rm -f client.csr client.crt client.key
rm -f ca.srl
# rm -f $(ls | egrep -v '.pem$')

if [ "$CA" != "" ]; then
    rm -f ca.pem
fi

cd $OPWD
ls -l $CERTS_DIR
