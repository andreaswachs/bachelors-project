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


