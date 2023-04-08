#!/usr/bin/env bash

baseurl='localhost:8008'

TOKEN=$(curl -s -d '{"username":"foo", "password":"bar"}' -H 'Content-Type: application/json' -X POST ${baseurl}/login | jq ".token" | sed 's/"//g')

curl -L -s --header "token: $TOKEN" ${baseurl}/deploy/ | jq

