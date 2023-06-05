#!/bin/bash

IP_PREFIX=172.16.86
DNS_IP=$IP_PREFIX.3

# Create the docker network
docker network create -d macvlan \
  --subnet "172.16.86.0/24" \
  --label experiment \
  isolated

## Note on network create:
## - removed gateway setting:
##   --gateway=172.16.86.1


# Start the DHCP service
docker run -dit --rm --init \
  --net isolated \
  --label experiment \
  --name isolated_dhcp \
  -v "$(pwd)/data":/data \
  networkboot/dhcpd:1.2.0 eth0

docker run -dit --rm --init \
  --net isolated \
  --label experiment \
  --mount type=bind,source=$(pwd)/Corefile,target=/Corefile \
  --mount type=bind,source=$(pwd)/zonefile,target=/zonefile \
  --name isolated_dns \
  --ip $DNS_IP \
  coredns/coredns:1.6.1

docker run -dit --rm --init \
  --net isolated \
  --label experiment \
  --name isolated_service_1 \
  --ip ${IP_PREFIX}.30 \
  --dns ${DNS_IP} \
  andreaswachs/placeholder_vuln_server:latest

docker run -dit --rm --init \
  --net isolated \
  --label experiment \
  --name isolated_actor \
  --ip ${IP_PREFIX}.31 \
  --dns $DNS_IP \
  praqma/network-multitool ash


