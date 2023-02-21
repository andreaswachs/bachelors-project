#!/bin/bash

# Shutting down containers
echo "Shutting down containers"
docker ps --filter "label=experiment" | tail -n +2 | awk '{ print $1}' | xargs -n1 docker kill

echo "Removing isolated network(s)"
docker network ls --filter "label=experiment" | tail -n +2 | awk '{ print $1}' | xargs -n1 docker network rm
