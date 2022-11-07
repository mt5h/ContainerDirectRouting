# How to build
```
chmod +x build.sh
./build.sh
```
# How to run

```
docker run --rm -p 8008:8008 -v /var/run/docker.sock:/var/run/docker.sock  cntspawner
```

# How it works
The cntSpawer is a simple Docker APIs wrapper to create, list and delete containers.

# Create a container

```
curl  --header "Content-Type: application/json" \
      --request POST \
      --data \
      '{ \
      "name":"test-5", \
      "network":"example", \
      "image":"app-instance", \
      "labels": \
      { \
      "traefik.enable": "true", \
      "traefik.http.routers.test-5.entrypoints": "web", \
      "traefik.http.routers.test-5.rule":"Path(`/test-5`)" \ 
      }, \
      "envs": {"INSTANCE":"test-5"}\
      }' \
localhost:8008/v1/

```

# Delete a container

```
curl -XDELETE localhost:8008/v1/$container_id
```
# List containers

```
curl localhost:8008/v1/
```
