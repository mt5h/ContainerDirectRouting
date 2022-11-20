#!/usr/bin/env bash

baseurl='localhost:8008/deploy'

response="$(curl -s ${baseurl}/ )"

if command -v jq; then
  echo $response | jq
else
  echo "Install jq for a better output"
  echo $response
fi
