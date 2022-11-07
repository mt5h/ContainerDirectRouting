#!/usr/bin/env bash

response="$(curl -s localhost:8008/v1/)"

if command -v jq; then
  echo $response | jq
else
  echo "Install jq for a better output"
  echo $response
fi
