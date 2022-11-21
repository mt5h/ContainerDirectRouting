#!/usr/bin/env bash

baseurl='localhost:8008/deploy'
instance_name='nightscout'
max=5

for i in $(seq 1 $max); do
  response=$(curl -L -s \
    --header "Content-Type: application/json" \
    -X DELETE \
    ${baseurl}/${instance_name}-$i
  )

  if command -v jq >/dev/null 2>&1; then
    echo $response | jq
  else
    echo "Install jq for a better output"
    echo $response
  fi
done
