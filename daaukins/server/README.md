# Daaukins Server

## Setup on hosts

Complete server setup and install: to setup dependencies, a daaukins user and downloading source files execute the following command as the root user on a freshly installed Ubuntu server:

```sh
curl -sL https://raw.githubusercontent.com/andreaswachs/bachelors-project/main/daaukins/server/initial_setup.sh | sh
```

You should find a `server` folder in user `daaukins` home folder.
This is a symbolic link to the embedded source folder in the repository `daaukins/server`.

You can configure the server with the `server.yaml` file.
Servers can be in "leader" or "follower" mode. The leader server has many followers, the followers can't have no followers.
The leader must have all followers configured.

## Running in dev mode

Running in a non-production mode is the only supported mode of running the Daaukins server.

Ensure that all Docker images are pulled down to the host before running:

```sh
make pull-images
```

Run the following command to start the service after you've properly configured `server.yaml` and `store.yaml`

```sh
make dev
```

## Interesting technical ~~limitations~~details

- Hosts should not have more than ~2TB of RAM, as the available memory calculations will break since they use ints to read the available memory in kB.
- Docker images needs to be pulled to the host machine before running.
