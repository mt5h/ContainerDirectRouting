## Architecture

This app spawn multiple containers of the minion app reachable directly trought the traefik proxy.

In order to to that the minion app accept requests for a specific path with the container name in it (but can be something else), like this:
```
/session/$container-name/ping
```

The traefik route will bind to the minon container with the rule PathPrefix(`/session/$container-name/`)

If the container does not exist it needs to be created by making a POST to /v1/ (see the example)

```
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
```

The minion app shutdown by itself after an IDLE time (default 1m). Can be ovveriden with the env var IDLE.

If the minon app exists but is shutdown, making a request to /session/$container-name/ping will bring it up, because the spawner app has a routing to /session/:container-name/.
This works because no specific route exist in traefik for the stopped container.

When the stopped container is up again a new (and specific) route is added automatically by traefik [see](https://doc.traefik.io/traefik/routing/routers/#rule), which allow direct routing.

```

---- request ---> traefik -- /v1/ --> spawner ----> docker ---->  minion-1
                     |                                     ---->  minion-2
                     |                                     ---->  minion-3
                     |
                     +-------/session/minon-$i/ping ----------->

```

## How it works

Build minionApp and cntSpawner 

```
cd minionApp
./build.sh
cd ..
cd cntSpawner
./build.sh
cd ..
```

run the compose file.

```
docker-compose up -d
```

This command creates the basic containers (traefik and spawner) in the specified attachable network.

## Start the minions

```
./create-minions.sh
```

Check the routes on traefik with the dashbord at localhost:8080

## Test your minions routing

```
./test-minions.sh
```

Wait for a minute until they shutdown.

(check with docker ps or the traefik dashbord)

Open your browser the go to:

```
http://localhost/session/minion-3/ping 
```

The request hang for a little and the redirects you to the service.

The following requests will work as usual.

## Clean up

```
./delete-minons.sh
```

terminate the docker compose

```
docker-compose down
```

## Security notice 
This software is a PoC not good for a production use, for example any stopped container on your machine can be started with this API.
The /v1/ should be protected from external calls.








