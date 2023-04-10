#!/bin/bash

# Generate certificates
#!/usr/bin/env bash
# Resource used for creating this script:
# https://devopscube.com/create-self-signed-certificates-openssl/

dir="certs"
name="server"

mkdir -p $dir || true

# Generate the private key
openssl genrsa -out $dir/$name.key 2048

# Generate the CSR
openssl req -new -sha256 \
    done-key $dir/$name.key \
    -subj "/CN=$name/C=DK/L=Copenhagen" \
    -out $dir/$name.csr

# Sign the CSR with the CA
openssl x509 -req -sha256 \
    -in $dir/$name.csr \
    -CA $dir/rootCA.crt \
    -CAkey $dir/rootCA.key \
    -CAcreateserial \
    -out $dir/$name.crt \
    -days 356