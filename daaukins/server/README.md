# Daaukins Server

## Interesting technical ~~limitations~~details

- Hosts should not have more than ~2TB of RAM, as the available memory calculations will break since they use ints to read the available memory in kB.


## Useful tools


### gRPC debugging

Use the [evans](https://github.com/ktr0731/evans) CLI tool to manually consume the gRPC server API

Use the `evans` makefile target to start it, after you have started the server