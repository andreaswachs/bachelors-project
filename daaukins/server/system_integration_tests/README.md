# `docker-compose/`

This directory contains docker-compose files that are used to test different scenarios of servers running


## `leader-leader.yaml`

This case checks that two instances of servers in leader mode will complain when `leader-1` tries to connect to `leader-2`.

Look out for the error log message!