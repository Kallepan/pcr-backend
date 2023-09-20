#!/bin/bash
# This script launches the database for the development environment

docker-compose -f docker-compose.db.yaml down -v
docker-compose -f docker-compose.db.yaml up