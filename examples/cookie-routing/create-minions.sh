#!/usr/bin/env bash

baseurl='localhost:8008/deploy'
instance_name='minion'
image='mock-app:latest'

max=2

generate_post_data()
{
cat<<EOF
{
  "name":"$instance_name-$1",
  "network":"traefiknet",
  "image":"$image",
  "labels": {
    "health-check": "http:\/\/$instance_name-$i:9000\/status",
    "traefik.enable": "true",
    "traefik.http.services.$instance_name-$1.loadbalancer.server.port": "9000",
    "traefik.http.routers.$instance_name-$1.entrypoints": "websecure",
    "traefik.http.routers.$instance_name-$1.rule": "HeadersRegexp(\"Cookie\",\".*instance=$instance_name-$1.*\")",
    "traefik.http.routers.$instance_name-$1.tls": "true"
  },
  "envs": {
			"TZ": "Etc\/UTC",
      "CONNSTR":"db $1",
      "IDLE":"1m"
	}
}
EOF
}


for i in $(seq 1 $max); do 
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
