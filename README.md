# Shortly
___

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

## Code Quality
#### Run lint
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

### Swagger

Install swagger
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```
Swagger Initialize
```shell
swag init
```
#### API DOC
First of all you must run the project.
http://localhost:6161/api/swagger/index.html

### Inception
[excalidraw](https://excalidraw.com/#json=gww-IAkHNXslZjEIwC4US,YrQX5EU2s--dIM9eeOrIBA)

### Used Libraries

[Swagger](https://github.com/swaggo) For API Documentation </br>
[Testify](https://github.com/stretchr/testify) For Test
