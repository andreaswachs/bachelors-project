# Daaukins client

The Daaukins client is a CLI tool that lets administrators interact with the Daaukins server(s).

This tool allows for administrators to perform CRUD operations on labs.

## Building the binary

To simply build the `dkn` binary and leave it in the `daaukins/client` directory

```sh
make build
```

## Installing

If you want the `dkn` binary installed to your system run:

```sh
make install
```

This moves the binary to `~/opt/bin` by default. You can set the variable `DKN_INSTALL_DIR` to change this when running the command:

```sh
DKN_INSTALL_DIR="~/my/other/bins" make install
```

## Actions

Here, I will showcase some useful and possible commands:

- `dkn get labs`: shows running labs
- `dkn get $ID`: shows information about a running lab given its id
- `dkn create lab -f filename.yaml`: creates a lab from a configuration file
- `dkn create lab -f -`: Reads from stdin for a passed configuration
- `dkn remove lab $id`: removes a running lab given its id
- `dkn config show`: shows the configured server that the CLI tools connects to
