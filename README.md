## Architecture

This app span multiple containers of the mock-app reachable directly trought the traefik proxy.
In order to do that the mock-app accept requests for a specific path with the container name in it (but can be something else), like this:

```
/session/$container-name
```

The traefik route will bind to the mock-app container with the rule PathPrefix(`/session/$container-name/`)
If the container does not exist it needs to be created by making a POST to the /deploy API (see the example)

```
  {
  "name":"minion-$1",
  "network":"traefiknet",
  "image":"minionapp:latest",
  "labels": {
    "healthcheck": "http:\/\/minion-$1:9000\/status",
    "traefik.enable": "true",
    "traefik.http.routers.minion-$1.entrypoints": "web",
    "traefik.http.routers.minion-$1.rule":"PathPrefix(\"/session/minion-$1/\")"
  },
  "envs": {"CONNSTR":"db $1", "IDLE":"1m"}
  }
```

If you specify the label 

```
"healthcheck": "http:\/\/minion-$1:9000\/status",
```

The spawner app will use the value of the label to perform an health check when the app is resumed. This help us to redirect the user at the proper time.

The mock-app shutdown by itself after an IDLE time (default 1m). Can be ovveriden with the env var IDLE.

If the mock-app exists but is shutdown, making a request to /session/$container-name will bring it up again, because the spawner app will answer to the route /session/:container-name/.
This works because no specific route exist in traefik for the stopped container.

When the stopped container is up again a new (and specific) route is added automatically by traefik [see](https://doc.traefik.io/traefik/routing/routers/#rule), which allow direct routing.

```

---- request ---> traefik -- /deploy/ --> spawner ----> docker ---->  minion-1
                     |                                         ---->  minion-2
                     |                                         ---->  minion-3
                     |
                     +-------/session/minon-$i ----------->

```

## How it works

Build mock-app and spawner containers

clone the git repo

```
cd ContainerDirectRouting
./build_project.sh
```

run the docker compose file.

```
docker-compose up -d
```

This command creates the basic containers (traefik and spawner) in the specified attachable network.

## Create some mock-app instances

Move to the examples folder

```
cd examples
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
http://localhost/session/minion-3
```
or launch again

```
./test-minions.sh
```
The request will hang for a little and then the spawner will redirects you to the correct container. All the other request will work normally.
The health-check internally checks the http://minion-x:9000/status API using the container networking and name resolution system, no ports of the mock-app need to be available outside the docker net.
The only constraint is that you should put the spawner and the mock-app on the same docker network.
If no health-check label is specified during the mock-app instance creation the redirect will wait for 2 seconds (by default), and the HTTP probe will be skipped.
The only endpoint reachable from outside the container network should be the ones mapped in traefik.
The spawner app /deploy API should be private. The /session/$container-name can be called from the outside to restore the container instance.
If you want to write your application that intercept the /session/$container-name request you can make PUT to /deploy/$container-ID to ask the spawner to restart the instace for you (no auto redirect).

## Clean up

```
./delete-minons.sh
cd ..
```

terminate the docker compose

```
docker-compose down
```

## Security notice 
This software is a PoC, APIs are not authenticated and should be managed carefully.
The /deploy/ route should be protected from external calls.
The /session/$container-name route is for public use.
This separation is done by traefik and is configured into the docker-compose.








