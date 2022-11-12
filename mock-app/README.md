# How to build

chmod +x build.sh
./build.sh

# How to run
docker run --rm -p 9000:9000 -e CONNSTR=FOO minionapp

# How it works
The minionapp listen on the port 9000 for HTTP requests.
It answer only to a specific path:

``` 
/$hostname/ping
```

The application routing is being built dinamically from the Gin Framework based on the hostname.

It returns the content of the CONNSTR variable.
