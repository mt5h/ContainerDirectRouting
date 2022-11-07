# How it works

Build minionApp and cntSpawner 

```
cd minionApp
./build.sh
cd ..
cd cntSpawner
./build.sh
cd ..
```

run the compose file

```
docker-compose up -d
```

This command creates the basic containers (traefik and spawner) in the specified attachable network

Start the your minions

```
./create-minions.sh
```

Check the routes on traefik with the dashbord at localhost:8080

Test you minions routing

```
./test-minions.sh
```

Clean up

```
./delete-minons.sh
```

terminate the docker compose

```
docker-compose down
```


