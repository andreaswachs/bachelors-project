# Experiment: isolated network

This folder contains files for an experiment where I manually mock a network as Haaukins creates it.

This experiment creates a virtual `macvlan` network, in which containers are assigned MAC addresses to make them look
like physical devices on the network.

The network is isolated in that there are no access between network barriers in either way.

The same DHCP and DNS services that Haaukins use are deployed with minimal configuration to have a custom domain `semifreddo.yum` resolve to
a docker container named `isolated_service_1`.

`isolated_service_1` is an instance of Aline Linux that is deployed with no custom configuration other than it has a specific IP set.

A container with networking tools are also deployed named `isolated_actor`.

The `isolated_service_1` and `isolated_actor` containers are set to do DNS resolution agains the custom DNS service.

## Running the experiment

### Setting up the experiment:

To setup the network and containers:

```sh
make s
```

or

```sh
make setup
```

### Attaching to the `isolated_actor`

In order to verify that the custom domain resolves and that the isolated network doesn't have access to the public internet you need to attach to an instance of a shell on the `isolated_actor` container:

```sh
docker attach isolated_actor
```

Here you can then ping the `isolated_service_1` container:

```sh
/ # ping -c 2 semifreddo.yum
PING semifreddo.yum (172.16.86.30) 56(84) bytes of data.
64 bytes from isolated_service_1.isolated (172.16.86.30): icmp_seq=1 ttl=64 time=0.071 ms
64 bytes from isolated_service_1.isolated (172.16.86.30): icmp_seq=2 ttl=64 time=0.063 ms

--- semifreddo.yum ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1020ms
rtt min/avg/max/mdev = 0.063/0.067/0.071/0.004 ms
```

### Restarting the experiment

If you've made changes and need to restart the experiment over again you can use:

```sh
make r
```

or

```sh
make refresh
```

### Cleaning up

When done experimenting, you can shut down the docker containers and remove the docker network with:

```sh
make c
```

or 

```sh
make cleanup
```

## Notes

### DHCP server

I am not sure that the DHCP server is doing anything as I am manually assigning IP addresses. Maybe it comes into play with the virtualbox VMs?

### IP generation

Haaukins keeps track of and generates IPs for these isolated networks. 

### Static IPs

It seems that the DNS service always are allocated with ip XXX.XXX.XXX.3


