#!/bin/bash

echo "Shutting down containers"
docker ps --filter "label=daaukins" | tail -n +2 | awk '{ print $1}' | xargs -r -n1 -P4 docker kill
# >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> [1] [2] [3]
# [1] -r: xargs doesn't execute if nothing was piped to it
# [2] -n1: execute the command once for each line
# [3] -P4: execute the command in parallel with 4 processes

echo "Removing isolated network(s)"
docker network ls --filter "label=daaukins" | tail -n +2 | awk '{ print $1}' | xargs -r -n1 -P4 docker network rm 