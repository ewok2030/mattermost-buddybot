# build stage
FROM golang:alpine AS build-env

# need git to run 'go get'
RUN apk add --no-cache git

WORKDIR /go/src
COPY . .
RUN go get -d -v ./...
RUN go build -v -o main

# final stage
FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/main /app/
ENTRYPOINT ./main