#!/usr/bin/env bash

which git > /dev/null 2>&1

if [[ $? -ne  0 ]]; then 
  echo "git is not installed using default tag"
  docker build -t spawner:latest -f ./spawner/Dockerfile .
  docker build -t mock-app:latest -f ./mock-app/Dockerfile .
else
  tag=${VERSION:-$( git describe --tags --dirty --abbrev=14 | sed -E 's/-([0-9]+)-g/.\1+/' )}
  docker build -t spawner:$tag -f ./spawner/Dockerfile .
  docker build -t mock-app:$tag -f ./mock-app/Dockerfile .

  docker tag spawner:$tag spawner:latest
  docker tag mock-app:$tag mock-app:latest
fi



