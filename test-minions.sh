#!/usr/bin/env bash

for i in {1..5}; do 
    response=$(curl -s --header "Content-Type: application/json" \
      localhost/session/minion-$i/ping)

  if command -v jq > /dev/null 2>&1; then
    echo $response | jq
  else
    echo "Install jq for a better output"
    echo $response
  fi
done
