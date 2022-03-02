# Shortly

## Setup

## Run lint
```shell
golangci-lint run -c .golangci.yml -v
```

## Dockerize
#### Build
```shell
docker build -t shortly .
```
#### Run
```shell
docker run -p 6161:6161 shortly
```
