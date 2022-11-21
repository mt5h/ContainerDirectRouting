#!/usr/bin/env bash

baseurl='localhost:8008/deploy'
instance_name='nightscout'
max=5

generate_post_data()
{
cat<<EOF
{
  "name":"$instance_name-$1",
  "network":"traefiknet",
  "image":"nightscout/cgm-remote-monitor:latest",
  "labels": {
    "health-check": "http:\/\/$instance_name-$i",
    "traefik.enable": "true",
    "traefik.http.services.$instance_name-$1.loadbalancer.server.port": "80",
    "traefik.http.routers.$instance_name-$1.entrypoints": "web",
    "traefik.http.routers.$instance_name-$1.rule": "HeadersRegexp(\"Cookie\",\".*instance=$instance_name-$1.*\")"
  },
  "envs": {
      "CUSTOM_TITLE": "nightscout-$1",
			"INSECURE_USE_HTTP": "true",
			"NODE_ENV": "production",
			"PORT": "80",
			"TZ": "Etc\/UTC",
			"MONGODB_URI": "mongodb:\/\/mongo:27017\/nightscout-$1",
			"API_SECRET": "change_me_please-$1",
			"ENABLE": "careportal rawbg iob",
			"AUTH_DEFAULT_ROLES": "readable"	
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
