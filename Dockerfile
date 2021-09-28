FROM golang:1.17 as base
RUN groupadd --gid 1001 detecc \
    && useradd --uid 1001 --gid detecc --shell /bin/bash --create-home detecc

WORKDIR /detecctor/src
RUN mkdir "/detecctor/plugins" && chown -R detecc:detecc /detecctor/plugins
USER detecc
ENV GOPATH /home/detecc/go
ENV GOBIN /home/detecc/go/bin
ENV GOCACHE /home/detecc/.cache
VOLUME /home/detecc/.cache
VOLUME /home/detecc/go
COPY . .

FROM base as dev
CMD ["go", "run", "."]

FROM base as run

ARG PLUGIN_DIR
#ENV PLUGIN_DIR=${PLUGIN_DIR}
RUN echo $PLUGIN_DIR
COPY $PLUGIN_DIR ../plugins

RUN go build main.go -o detecctor

FROM base as test
RUN go test -v

FROM alpine as app

COPY --from=run /detecctor /detecctor
RUN mv /detecctor/src/detecctor /usr/bin/detecctor
ENTRYPOINT ["detecctor"]