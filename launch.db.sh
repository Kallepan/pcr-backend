#!/bin/bash

cd psql
docker-compose down -v
docker-compose up -d --build
