#!/usr/bin/env bash
# Check versions
docker version
docker-compose version
git version

# Checkout the branch
rm -rf discount
git clone https://github.com/jsotogaviard/discount
cd discount

# Start docker
docker-compose up -d --build