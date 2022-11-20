#!/usr/bin/env bash

baseurl='localhost:8008/deploy'

for i in {1..5}; do 
  response=$(curl -L -s --header "Content-Type: application/json" \
     -X DELETE \
     ${baseurl}/minion-$i
   )

  if command -v jq >/dev/null 2>&1; then
    echo $response | jq
  else
    echo "Install jq for a better output"
    echo $response
  fi
done
