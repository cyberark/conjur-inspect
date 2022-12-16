#!/usr/bin/env bash

# Start the docker daemon in the containr
service docker start

exec "$@"
