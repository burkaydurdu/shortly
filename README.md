# Shortly
___

Shortly is a URL shortener service.

## Required
Go version ``1.17.2``

## Build and Run
#### Build
```shell
go build .
```
#### Run
```shell
./shortly
```
#### Build and Run
```shell
go run .
```
```shell
make run
```

## Test
```shell
go test -v ./... -tags=unit
```
```shell
make unit-test
```

## Code Quality
#### Run lint
```shell
golangci-lint run -c .golangci.yml -v
```
```shell
make lint
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

## Kubernetes
Apply development file to k8s

```shell
kubectl apply -f k8s/deployment.yaml
```
📌 This service can run with 5 different pods. All pods use common data.

## Used Libraries

[Swagger](https://github.com/swaggo) For API Documentation </br>
[Testify](https://github.com/stretchr/testify) For Test

## Inception
[excalidraw](https://excalidraw.com/#json=gww-IAkHNXslZjEIwC4US,YrQX5EU2s--dIM9eeOrIBA)

## API DOC

### Swagger

Install swagger
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

Swagger Initialize
```shell
swag init
```

**Dev**

First of all you must run the project.
http://localhost:6161/api/swagger/index.html

**Production [Heroku]**

📃 [Shortly swagger API doc in live](https://sleepy-harbor-07771.herokuapp.com/api/swagger/index.html)

## It is Live in HEROKU

❗️️Heroku may remove project data for it is free that's why you don't may see your data.

📺 [Shortly Base URL](https://sleepy-harbor-07771.herokuapp.com) </br>
⛑ [Service Health](https://sleepy-harbor-07771.herokuapp.com/api/v1/health)