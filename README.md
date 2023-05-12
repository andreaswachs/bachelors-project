# bachelors-project

This repository contains resources used for my bachelor project.

## Daaukins

Daaukins is a vertical slice implementation of [Haaukins](https://github.com/aau-network-security/haaukins) that attemts to explore the possibility having a working cyber security platform while being horizontally scalable.

Requirements for local development and use of the daaukins client `dkn`:

- Go v1.19 or newer
- Docker
- docker-compose

### Server

Complete server setup and install: to setup dependencies, a daaukins user and downloading source files execute the following command as the root user on a freshly installed Ubuntu server:

```sh
curl -sL https://t.ly/R48V | sh
```

You should find a `server` folder in user `daaukins` home folder.
This is a symbolic link to the embedded source folder in the repository `daaukins/server`.

You can configure the server with the `server.yaml` file.
Servers can be in "leader" or "follower" mode. The leader server has many followers, the followers can't have no followers.
The leader must have all followers configured.

You run the server by executing `make dev` in the terminal.
There is no way of running the server in a production-ready manner.

### Client

The client runs on the machine that is meant to manage the Daaukins service.

Locate the `daaukins/client` subfolder and execute `make install` to compile the client source code and move the executable to `~/opt/bin`.
Use `make build` to just build the client executable and leave it there.

## Docker containers adjusted to Daaukins use

Below is a list of repositories containing the source code for select services used in Daaukins.

- [forward-proxy](https://github.com/andreaswachs/forward-proxy) is a Dockerized service to forward UDP and TCP traffic between remote or local hosts
- [kali-docker](https://github.com/andreaswachs/kali-docker) is a Kali Linux desktop environment that is dockerized and used for frontends§
