#!/bin/bash
set -e

docker-compose down --remove-orphans
docker-compose rm -fv

rm -rf ./postgres-data
docker-compose up --build
