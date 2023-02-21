# Pros and Cons

## Nomad

Pros:
* Can orchestrate both containers and VMs

Cons:
* Unable to assign static IP addresses using a custom Container Network Interface (CNI) for working with `macvlan` networks
* Configuration requires Consul, another service working
* Can only orchestrate QEMU VMs with a third-party plugin. Customized networking is not yet supported

## Docker Swarm

Pros:
* ?

Cons:
* Uses overlay networks to allow for communication between swarm nodes. This might break low level exploits as traffic is routed over higher level transportation over the network
* Other than scalability and services, it doesn't really provide anything useful as there still are a lot of interaction with Docker Networks

## Kubernetes

Pros:
* Can control network access through network policy tools such as Calico

Cons:
* No fine grained control over subnets/ips. Can't be sure to support the foundational features of Haaukins of randomized subnets.

