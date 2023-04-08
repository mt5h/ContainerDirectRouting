#!/usr/bin/env bash

baseurl='localhost:8008/deploy'
instance_name='minion'
max=2


baseurl='localhost:8008'

TOKEN=$(curl -s -d '{"username":"foo", "password":"bar"}' -H 'Content-Type: application/json' -X POST ${baseurl}/login | jq ".token" | sed 's/"//g')

for i in $(seq 1 $max); do
  response=$(curl -L -s \
    -v \
    --header "token: $TOKEN" \
    --header "Content-Type: application/json" \
    -X DELETE \
    ${baseurl}/deploy/${instance_name}-$i
  )

done
