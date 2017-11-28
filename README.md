# GOTTER

This is an implementation of a microservice architecture following [this blog series](http://callistaenterprise.se/blogg/teknik/2017/02/17/go-blog-series-part1/)

## Running

1. Install [docker](https://docs.docker.com/engine/installation/) and [docker-compose](https://docs.docker.com/compose/install/).

2. Enable swarm mode

```sh
docker swarm init
```

3. Create an account in [Loggly](https://loggly.com) and get a Customer Token.

4. Prepare a go environment and get the code:

```sh
export GOPATH=<some_path>
mkdir -p $GOPATH/src/github.com/alculquicondor
cd $GOPATH/src/github.com/alculquicondor
git clone https://github.com/alculquicondor/gotter
cd gotter
```

5. Paste the Loggly token in `gelftail/token.txt`

6. Build

```
make
```

6. Run the stack in docker swarm

```sh
make deploy
```

7. Check the stack and or services

```sh
docker stack ps gotter
docker service ls -baseAddr=127.0.0.1 -zuul=false
```

8. Visualize containers in http://127.0.0.1:8080

## Load Testing

1. Build the load generator

```sh
make loadtest/loadtest
./loadtest/loadtest -
```
