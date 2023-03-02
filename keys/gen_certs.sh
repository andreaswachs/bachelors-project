#!/usr/bin/env bash
# Resource used for creating this script:
# https://devopscube.com/create-self-signed-certificates-openssl/

dir="keys"

# RootCA
## Clean up
rm -rf $dir/ca || true
mkdir -p $dir/ca

## Generate a CA certificate
openssl req -x509 \
            -sha256 -days 356 \
            -nodes \
            -newkey rsa:2048 \
            -subj "/CN=SEC1MA2/C=DK/L=Copenhagen" \
            -keyout $dir/ca/rootCA.key -out $dir/ca/rootCA.crt 


for host in "example.com"
do
    # Clean up
    rm -rf $dir/$host || true
    mkdir -p $dir/$host

    # Generate the private key
    openssl genrsa -out $dir/$host/$host.key 2048

    # Generate the CSR
    openssl req -new -sha256 \
                -key $dir/$host/$host.key \
                -subj "/CN=$host/C=DK/L=Copenhagen" \
                -out $dir/$host/$host.csr

    # Sign the CSR with the CA
    openssl x509 -req -sha256 \
                 -in $dir/$host/$host.csr \
                 -CA $dir/ca/rootCA.crt \
                 -CAkey $dir/ca/rootCA.key \
                 -CAcreateserial \
                 -out $dir/$host/$host.crt \
                 -days 356
done
