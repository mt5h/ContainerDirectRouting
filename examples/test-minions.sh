#!/usr/bin/env bash

baseurl='localhost'

for i in {1..5}; do 
    echo "testing: ${baseurl}/minion-$i/"
    response=$(curl -L -s --header "Content-Type: application/json" \
      "${baseurl}/minion-$i/"
    )

  if command -v jq > /dev/null 2>&1; then
    echo $response | jq
  else
    echo "Install jq for a better output"
    echo $response
  fi
done
