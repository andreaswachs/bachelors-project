# Daaukins Server

## Setup

TODO: write how to setup a server

``` sh
curl -sL https://raw.githubusercontent.com/andreaswachs/bachelors-project/feature/main/daaukins/server/initial_setup.sh | sh
```



## Interesting technical ~~limitations~~details

- Hosts should not have more than ~2TB of RAM, as the available memory calculations will break since they use ints to read the available memory in kB.

