# bachelors-project

This repository contains resources used for my bachelor project.

## Daaukins

Daaukins is a vertical slice implementation of [Haaukins](https://github.com/aau-network-security/haaukins) that attemts to explore the possibility having a working cyber security platform while being horizontally scalable.

The Daaukins server is only functional on linux due to linux only support for macvlan networks.

Requirements for local development and use of the daaukins client `dkn`:

- Go v1.19 or newer
- Docker
- docker-compose

### Server

See more in the `daaukins/server` directory

### Client

The client runs on the machine that is meant to manage the Daaukins service.

Locate the `daaukins/client` subfolder to see more.

## Auxiliary repositories with project dependencies

Below is a list of repositories containing the source code for select services used in Daaukins.

- [daaukins-service](https://github.com/andreaswachs/daaukins-service) is a repository containing the gRPC interface files that is used for both the Daaukins server and client. Golang package documentation can be found [here](https://pkg.go.dev/github.com/andreaswachs/daaukins-service)
- [forward-proxy](https://github.com/andreaswachs/forward-proxy) is a Dockerized service to forward UDP and TCP traffic between remote or local hosts
- [kali-docker](https://github.com/andreaswachs/kali-docker) is a Kali Linux desktop environment that is dockerized and used for frontends§
- [vuln-server](https://github.com/andreaswachs/vuln-server) is a placeholder for a vulnerable server. It hosts a directory listing over HTTP on port 80, with a file called "flag.txt" with a randomly generated "daaukins" flag.
