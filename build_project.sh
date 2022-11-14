#!/usr/bin/env bash


which git > /dev/null 2>&1

if [[ $? -ne  0 ]]; then 
  tag="latest"
  echo "Git is not installed using default tag: $tag"
else
  # try to get the tag from the current branch
  tag=$(git describe --tags --abbrev=0 --dirty)
  
  if [[ $? -ne 0 ]]; then
    # git descibe failed set latest tag
    echo "Can not get the tag from the current repo, setting to latest"
    tag="latest"
  else
    tag="$tag-alpine"
  fi
fi


docker build -t spawner:$tag -f ./spawner/Dockerfile .
docker build -t mock-app:$tag -f ./mock-app/Dockerfile .

# mark the last built as latest
if [ $tag != "latest" ]; then
  docker tag spawner:$tag spawner:latest
  docker tag mock-app:$tag mock-app:latest
fi
