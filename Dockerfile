FROM golang:latest as base
WORKDIR /detecctor/src
COPY . .
RUN mkdir "/detecctor/plugins"

FROM base as dev
ENTRYPOINT ["go run ."]

FROM base as run

ARG PLUGIN_DIR
ENV PLUGIN_DIR=$PLUGIN_DIR

COPY $PLUGIN_DIR ../plugins

RUN go build main.go -o detecctor
ENTRYPOINT ["./detecctor"]

FROM base as test
RUN go test -v

FROM alpine as app

COPY --from=run /detecctor /detecctor
RUN mv /detecctor/src/detecctor /usr/bin/detecctor
ENTRYPOINT ["detecctor"]