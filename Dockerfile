#FROM golang:1.17.2-alpine as builder
#WORKDIR /go/src/app
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
#
#FROM scratch
#COPY --from=builder /go/src/app/app /app
#ENTRYPOINT ["/app"]

FROM scratch
COPY shortly /app
ENTRYPOINT ["/app"]
