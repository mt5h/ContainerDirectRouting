# How to build
```
chmod +x build.sh
./build.sh
```
# How to run

```
docker run --rm -p 8008:8008 -v /var/run/docker.sock:/var/run/docker.sock spawner
```

# How it works
The spawer is a simple Docker APIs wrapper to create, list and delete containers.

# Create a container

```
baseurl='localhost:8008/deploy'

generate_post_data()
{
  cat<<EOF
  {
  "name":"minion-$1",
  "network":"traefiknet",
  "image":"mock-app:latest",
  "labels": {
    "healthcheck": "http:\/\/minion-$1:9000\/status",
    "traefik.enable": "true",
    "traefik.http.routers.minion-$1.entrypoints": "web",
    "traefik.http.routers.minion-$1.rule":"PathPrefix(\"/session/minion-$1/\")"
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
```

# Delete a container

```
curl -XDELETE localhost:8008/deploy/$container_id
```
# List containers

```
curl localhost:8008/deploy/
```
