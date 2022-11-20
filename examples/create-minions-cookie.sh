#!/usr/bin/env bash

# TODO
# evaluate the creation of a single middleware
# can we exclude mininion from the regex?

baseurl='localhost:8008/deploy'

generate_post_data()
{
  cat<<EOF
  {
  "name":"minion-$1",
  "network":"traefiknet",
  "image":"mock-app:latest",
  "labels": {
    "health-check": "http:\/\/minion-$1:9000\/status",
    "traefik.enable": "true",
    "traefik.http.routers.minion-$1.entrypoints": "web",
    "traefik.http.routers.minion-$1.rule": "HeadersRegexp(\"Cookie\",\".*instance=minion-$1.*\")"
  },
  "envs": {"CONNSTR":"db $1", "IDLE":"1m"}
  }
EOF
}


for i in {1..5}; do 
response=$(curl -L -s --header "Content-Type: application/json" \
     -X POST \
     --data  "$(generate_post_data $i)" \
     ${baseurl}/
   )

  if command -v jq > /dev/null 2>&1; then
    echo $response | jq
  else
    echo "Install jq for a better output"
    echo $response
  fi
done
