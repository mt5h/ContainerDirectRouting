#!/usr/bin/env bash


for i in {1..5}; do 
  response=$(curl -s --header "Content-Type: application/json" \
     -X DELETE \
     localhost:8008/v1/minion-$i)

  if command -v jq >/dev/null 2>&1; then
    echo $response | jq
  else
    echo "Install jq for a better output"
    echo $response
  fi
done
