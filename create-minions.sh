#!/usr/bin/env bash

generate_post_data()
{
  cat<<EOF
  {
  "name":"minion-$1",
  "network":"traefiknet",
  "image":"minionapp:latest",
  "labels": {
    "traefik.enable": "true",
    "traefik.http.routers.minion-$1.entrypoints": "web",
    "traefik.http.routers.minion-$1.rule":"PathPrefix(\"/session/minion-$1/\")"
  },
  "envs": {"CONNSTR":"db $1", "IDLE":"1m"}
  }
EOF
}



for i in {1..5}; do 
response=$(curl -s --header "Content-Type: application/json" \
     -X POST \
     --data  "$(generate_post_data $i)" localhost/v1/)

  if command -v jq > /dev/null 2>&1; then
    echo $response | jq
  else
    echo "Install jq for a better output"
    echo $response
  fi
done
