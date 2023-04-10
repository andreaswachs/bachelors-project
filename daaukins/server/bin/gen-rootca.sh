#!/bin/bash

dir="certs"

# RootCA
## Clean up
rm -rf $dir || true
mkdir -p $dir

## Generate a CA certificate
openssl req -x509 \
            -sha256 -days 356 \
            -nodes \
            -newkey rsa:2048 \
            -subj "/CN=DAAUKINS/C=DK/L=Copenhagen" \
            -keyout $dir/rootCA.key -out $dir/rootCA.crt 
